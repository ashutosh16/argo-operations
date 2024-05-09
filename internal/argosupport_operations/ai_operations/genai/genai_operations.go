package genai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/argoproj-labs/argo-operations/api/v1alpha1"
	"github.com/argoproj-labs/argo-operations/internal/argosupport_operations"
	"github.com/argoproj-labs/argo-operations/internal/utils"
	rolloutv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
)

const (
	appSecretKey         = "app.secret"
	genAIEndPointSuffix  = "/analyze"
	argoV1ResAPI         = "/api/v1/applications"
	responseType         = "/analyses"
	argocdEndPointSuffix = "/api/v1/applications/"
	rolloutRevision      = "rollout.argoproj.io/revision"
)

type GenAIOperation struct {
	k8sClient     client.Client
	genAIClient   HttpClient
	dynamicClient dynamic.DynamicClient
	argoCDClient  HttpClient
	kubeClient    kubernetes.Interface
}

var (
	_ argosupport_operations.Executor = &GenAIOperation{}
)

func NewGenAIOperations(ctx context.Context, k8sClient client.Client, dynamicClient dynamic.DynamicClient, kubeClient kubernetes.Interface, wf *v1alpha1.Workflow, obj metav1.Object) (*GenAIOperation, error) {
	//logger := log.FromContext(ctx)
	authProviders, err := utils.GetAIProviders(ctx, k8sClient, &wf.Ref, obj)
	if err != nil {
		return nil, err
	}

	genClient, err := getGenAIClientWithSecret(ctx, k8sClient, authProviders, obj)
	if err != nil {
		return nil, err
	}

	argoCDClient, err := getArgoCDClienWithSecret(ctx, k8sClient, authProviders, obj)
	if err != nil {
		return nil, err
	}

	return &GenAIOperation{
		k8sClient:     k8sClient,
		genAIClient:   *genClient,
		argoCDClient:  *argoCDClient,
		dynamicClient: dynamicClient,
		kubeClient:    kubeClient,
	}, nil
}

func getGenAIClientWithSecret(ctx context.Context, k8sClient client.Client, authProviders *[]v1alpha1.AuthProvider, obj metav1.Object) (*HttpClient, error) {
	logger := log.FromContext(ctx)
	for _, authProvider := range *authProviders {
		if authProvider.Name == "genai-auth-provider" {
			secret, err := utils.GetSecret(ctx, k8sClient, &authProvider)
			if err != nil {
				logger.Error(err, "failed to get Secret from AuthProvider", "namespace", obj.GetNamespace(), "name", authProvider.Name)
				return nil, err
			}

			if secret == nil {
				return nil, fmt.Errorf("secret is missing")
			}

			return &HttpClient{
				BaseURL:          authProvider.Spec.Auth.BaseURL,
				AppID:            authProvider.Spec.Auth.AppID,
				IdentityEndpoint: authProvider.Spec.Auth.IdentityEndpoint,
				IdentityJobID:    authProvider.Spec.Auth.IdentityJobID,
				APIVersion:       authProvider.Spec.Auth.APIVersion,
				AppSecret:        string(secret.Data[appSecretKey]),
			}, nil
		}
	}
	return nil, nil
}

func getArgoCDClienWithSecret(ctx context.Context, k8sClient client.Client, authProviders *[]v1alpha1.AuthProvider, obj metav1.Object) (*HttpClient, error) {
	logger := log.FromContext(ctx)
	for _, authProvider := range *authProviders {
		if authProvider.Name == "argocd-auth-provider" {
			secret, err := utils.GetSecret(ctx, k8sClient, &authProvider)
			if err != nil {
				logger.Error(err, "failed to get Secret from AuthProvider", "namespace", obj.GetNamespace(), "name", authProvider.Name)
				return nil, err
			}

			if secret == nil {
				return nil, fmt.Errorf("secret is missing")
			}

			return &HttpClient{
				BaseURL:   authProvider.Spec.Auth.BaseURL,
				AppSecret: string(secret.Data[appSecretKey]),
			}, nil
		}
	}
	return nil, nil
}
func (g *GenAIOperation) Process(ctx context.Context, obj metav1.Object) (*v1alpha1.ArgoAISupport, error) {

	labels := obj.GetLabels()
	fullUrl := fmt.Sprint(g.argoCDClient.BaseURL + argocdEndPointSuffix + labels["app.kubernetes.io/instance"])

	app, err := g.argoCDClient.GetRequest(fullUrl, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to process the request: %v", err)
	}

	t, _ := g.buildAITokens(ctx, app, obj)

	//if err != nil {
	//	return nil, fmt.Errorf("error generating tokens")
	//}
	//tokensForAI := "{\n          \"failures\": [\n            {\n              \"context\":  }\n    t      ]\n        }"
	failures := Failures{
		Failures: []Failure{
			{Context: t},
		},
	}

	// Marshal the struct to JSON
	tokens, _ := json.Marshal(failures)
	res, err := g.genAIClient.PostRequest(ctx, string(tokens), genAIEndPointSuffix)
	if err != nil {
		return nil, fmt.Errorf("failed to post request: %v", err)
	}

	summary, ok := res.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("type assertion to map[string]interface{} failed")
	}

	value, exists := summary["analyses"]
	if !exists {
		return nil, fmt.Errorf("key 'analyses' not found in the result")
	}

	analysesSlice, ok := value.([]interface{})
	if !ok {
		return nil, fmt.Errorf("type assertion for 'analyses' as []interface{} failed")
	}

	// Assuming that each element in analysesSlice is a map[string]interface{} that contains an "analysis" key
	var genSummary string
	for _, analysis := range analysesSlice {
		analysisMap, ok := analysis.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("type assertion for individual analysis failed")
		}
		genSummary, ok = analysisMap["analysis"].(string)
		if !ok {
			return nil, fmt.Errorf("type assertion for 'analysis' as string failed")
		}
		break
	}

	argoOpsobj, ok := obj.(*v1alpha1.ArgoAISupport)
	if !ok {
		return nil, fmt.Errorf("type assertion to *v1alpha1.ArgoAISupport failed")
	}
	now := metav1.Now()

	argoOpsobj.Status = v1alpha1.ArgoAISupportStatus{
		Results: append(argoOpsobj.Status.Results, v1alpha1.Result{
			Name: argoOpsobj.Spec.Workflows[0].Name,
			Summary: v1alpha1.Summary{
				MainSummary: genSummary,
			},
			FinishedAt: &now,
			Phase:      v1alpha1.ArgoSupportPhaseCompleted,
			Message:    "Gen AI request completed",
		}),
		Phase: v1alpha1.ArgoSupportPhaseCompleted,
	}

	argoOpsobj.Status.Phase = "completed"
	return argoOpsobj, nil
}

func (g *GenAIOperation) buildAITokens(ctx context.Context, app *Application, o metav1.Object) (string, error) {
	logger := log.FromContext(ctx)

	var builder strings.Builder
	builder.WriteString(utils.GetInlinePrompt("app-conditions", ""))
	if app.Status.Health.Status == HealthStatusHealthy {
		builder.WriteString(utils.GetInlinePrompt("app-conditions", ""))
		if len(app.Status.Conditions) > 0 {
			for _, condition := range app.Status.Conditions {
				builder.WriteString(fmt.Sprintf("Condition Message: %s, Status: %s, LastTransitionTime: %s\n", condition.Type, condition.Message, condition.LastTransitionTime))
			}
		}
	}

	// 1. Fetch argo app data and analysis the health and app conditions and return the token
	// 2. Fetch the rollout information and check the health
	for _, res := range app.Status.Resources {
		if res.Health != nil && res.Health.Status != HealthStatusHealthy {
			builder.WriteString(fmt.Sprintf("Resource Name: %s Resource Health: %s  and kubernetes Message: %s", res.Name, res.Health.Status, res.Health.Message))
		}
	}

	rolloutLister := rolloutListFromClient(g.dynamicClient)
	res, err := rolloutLister(o.GetNamespace(), metav1.ListOptions{})

	var podList []string
	var aRuns []*rolloutv1alpha1.AnalysisRun

	if err != nil {
		return "", err
	} else {
		for _, r := range res {
			builder.WriteString(utils.GetInlinePrompt("rollout", r.Name))
			if r.Status.Phase != rolloutv1alpha1.RolloutPhaseHealthy {
				if rollout, ok := utils.StripTheKeys(r).(*rolloutv1alpha1.Rollout); ok {
					pods, _ := getPodsWithLabel(g.k8sClient, r.Status.CurrentPodHash)
					podList = append(podList, pods...)
					analysisLister := analysisListFromClient(g.dynamicClient)
					aRuns, err = analysisLister(o.GetNamespace(), metav1.ListOptions{})
					for i, ar := range aRuns {
						if o.GetAnnotations()[rolloutRevision] == ar.Annotations[rolloutRevision] {
							aRuns = append(aRuns[:i], aRuns[i+1:]...)
						}
					}
					builder.WriteString(rollout.Status.String())
				}
			} else {
				logger.Info("Rollout seems to be healthy and should not be included in the genai analysis")
			}
			if aRuns != nil && len(aRuns) > 1 {
				builder.WriteString(utils.GetInlinePrompt("analysis-runs", ""))
				// Check the latest revision
				builder.WriteString(aRuns[0].Status.String())
			}
			if podList != nil && len(podList) > 1 {
				// it's okay to just check only one pod, since the error is common
				logs, err := getLogsForPod(podList[0], r.Namespace, g.kubeClient)

				if err != nil {
					if strings.Contains(err.Error(), "no error found in logs") {
						builder.WriteString(utils.GetInlinePrompt("no-pod-error-log", ""))
					} else if strings.Contains(err.Error(), "could not") {
						logger.Error(err, "failed to process the pod logs")
					} else {
						builder.WriteString(utils.GetInlinePrompt("pod", ""))
						builder.WriteString(logs)
					}
				}

			} else {
				builder.WriteString(utils.GetInlinePrompt("no-pod-log", ""))
			}

		}
		if len(res) > 1 {
			builder.WriteString(utils.GetInlinePrompt("multi-rollout", ""))
		}
	}

	logger.Info("start collecting pod data")
	var eventList v1.EventList
	err = g.k8sClient.List(ctx, &eventList, client.InNamespace(o.GetNamespace()))

	if err != nil {
		logger.Error(err, "Failed to fetch events for namespace %s", o.GetNamespace())
	} else {
		for _, event := range eventList.Items {
			if event.Message == "Warning" || strings.Contains(event.Message, "Failed") {
				logger.Info("Failed Event Detected:", "Reason", event.Reason, "Message", event.Message)
				builder.WriteString(event.String())
			}
		}
	}

	return builder.String(), nil
}

func genericListFromClient(c dynamic.DynamicClient, gvr schema.GroupVersionResource) func(string, metav1.ListOptions) ([]*unstructured.Unstructured, error) {
	return func(namespace string, options metav1.ListOptions) ([]*unstructured.Unstructured, error) {
		res, err := c.Resource(gvr).Namespace(namespace).List(context.Background(), options)
		if err != nil {
			return nil, err
		}
		var resourceList []*unstructured.Unstructured
		for i := range res.Items {
			resourceList = append(resourceList, &res.Items[i])
		}
		return resourceList, nil
	}
}

type rolloutListFunc func(namespace string, options metav1.ListOptions) ([]*rolloutv1alpha1.Rollout, error)
type analysisListFunc func(namespace string, options metav1.ListOptions) ([]*rolloutv1alpha1.AnalysisRun, error)

func rolloutListFromClient(c dynamic.DynamicClient) rolloutListFunc {
	genericLister := genericListFromClient(c, v1alpha1.SchemeGroupVersion.WithResource("rollouts"))
	return func(namespace string, options metav1.ListOptions) ([]*rolloutv1alpha1.Rollout, error) {
		unstructuredList, err := genericLister(namespace, options)
		if err != nil {
			return nil, err
		}
		var rolloutList []*rolloutv1alpha1.Rollout
		for _, unstructuredRollout := range unstructuredList {
			rollout := &rolloutv1alpha1.Rollout{}
			err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredRollout.Object, rollout)
			if err != nil {
				return nil, err
			}
			rolloutList = append(rolloutList, rollout)
		}
		return rolloutList, nil
	}
}

func analysisListFromClient(c dynamic.DynamicClient) analysisListFunc {
	genericLister := genericListFromClient(c, v1alpha1.SchemeGroupVersion.WithResource("analysisruns"))
	return func(namespace string, options metav1.ListOptions) ([]*rolloutv1alpha1.AnalysisRun, error) {
		unstructuredList, err := genericLister(namespace, options)
		if err != nil {
			return nil, err
		}
		var analysisRunList []*rolloutv1alpha1.AnalysisRun
		for _, unstructuredAnalysisRun := range unstructuredList {
			analysisRun := &rolloutv1alpha1.AnalysisRun{}
			err = runtime.DefaultUnstructuredConverter.FromUnstructured(unstructuredAnalysisRun.Object, analysisRun)
			if err != nil {
				return nil, err
			}
			analysisRunList = append(analysisRunList, analysisRun)
		}
		return analysisRunList, nil
	}
}

func getPodsWithLabel(K8sClient client.Client, label string) ([]string, error) {
	podList := &v1.PodList{}

	// Create a ListOption with LabelSelector
	listOpts := []client.ListOption{
		client.MatchingLabels{"rollouts-pod-template-hash": label},
	}
	if err := K8sClient.List(context.TODO(), podList, listOpts...); err != nil {
		return nil, err
	}
	var podNames []string
	for _, pod := range podList.Items {
		podNames = append(podNames, pod.Name)
	}
	return podNames, nil
}
func getLogsForPod(podName, namespace string, kubeClient kubernetes.Interface) (string, error) {
	// Fetch logs from the pod
	podLogOpts := v1.PodLogOptions{}
	req := kubeClient.CoreV1().Pods(namespace).GetLogs(podName, &podLogOpts)

	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return "", fmt.Errorf("could not fetch logs: %v", err)
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	if _, err = io.Copy(buf, podLogs); err != nil {
		return "", fmt.Errorf("could not read logs: %v", err)
	}

	logs := buf.String()
	lines := strings.Split(logs, "\n")
	for i, line := range lines {
		if strings.Contains(line, "error") {
			start := maxLine(0, i-5)
			end := minLine(len(lines), i+6)
			return strings.Join(lines[start:end], "\n"), nil
		}
	}

	return "", fmt.Errorf("no error found in logs")
}

func minLine(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxLine(a, b int) int {
	if a > b {
		return a
	}
	return b
}
