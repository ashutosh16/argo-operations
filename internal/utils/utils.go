package utils

import (
	"context"
	"fmt"
	v1alpha1 "github.com/argoproj-labs/argo-support/api/v1alpha1"
	rolloutv1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var authProviderMap = make(map[v1alpha1.NamespacedObjectReference]*v1alpha1.AuthProvider)

func GetSecret(ctx context.Context, k8sClient client.Client, authProvider *v1alpha1.AuthProvider) (*v1.Secret, error) {
	logger := log.FromContext(ctx)

	var secret v1.Secret
	objectKey := client.ObjectKey{
		Namespace: authProvider.GetNamespace(),
		Name:      authProvider.Spec.SecretRef.Name,
	}
	err := k8sClient.Get(ctx, objectKey, &secret)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Secret for AuthProvider not found", "namespace", objectKey.Namespace, "name", objectKey.Name)
			return nil, err
		}

		logger.Error(err, "failed to get Secret from AuthProvider", "namespace", objectKey.Namespace, "name", objectKey.Name)
		return nil, err
	}

	return &secret, nil
}

func getAuthProviders(ctx context.Context, k8sClient client.Client, labels map[string]string) (*[]v1alpha1.AuthProvider, error) {
	logger := log.FromContext(ctx)

	authProviderList := &v1alpha1.AuthProviderList{}
	listOptions := []client.ListOption{
		client.MatchingLabels(labels),
	}
	err := k8sClient.List(ctx, authProviderList, listOptions...)
	if err != nil {
		return nil, err
	}

	if len(authProviderList.Items) == 0 {
		err = fmt.Errorf("authProvider not found with labels %#v", labels)
		logger.Info(err.Error())
		return nil, err
	}

	return &authProviderList.Items, nil
}

func GetAuthProviders(ctx context.Context, k8sClient client.Client, refs *[]v1alpha1.NamespacedObjectReference, namespace string) (*[]v1alpha1.AuthProvider, error) {
	logger := log.FromContext(ctx)
	authProviders, err := getAuthProviders(ctx, k8sClient, map[string]string{v1alpha1.LabelKeyAppName: v1alpha1.LabelKeyAppNameValue})
	if err != nil {
		logger.Error(err, "failed to get AuthProvider", "namespace", namespace)
		return nil, err
	}
	return authProviders, nil
}

func getConfigMap(ctx context.Context, k8sClient client.Client, labels map[string]string) (*v1.ConfigMap, error) {
	logger := log.FromContext(ctx)

	cmList := &v1.ConfigMapList{}
	listOptions := []client.ListOption{
		client.MatchingLabels(labels),
	}
	err := k8sClient.List(ctx, cmList, listOptions...)
	if err != nil {
		return nil, err
	}

	if len(cmList.Items) == 0 {
		err = fmt.Errorf("ConfigMap not found with labels %#v", labels)
		logger.Info(err.Error())
		return nil, err
	}

	return &cmList.Items[0], nil
}

func GetConfigMapRef(ctx context.Context, k8sClient client.Client, refs *v1alpha1.ConfigMapRef, namespace string) (*v1.ConfigMap, error) {
	logger := log.FromContext(ctx)
	cm, err := getConfigMap(ctx, k8sClient, map[string]string{v1alpha1.LabelKeyAppName: v1alpha1.LabelKeyAppNameValue})
	if err != nil {
		logger.Error(err, "failed to get AuthProvider", "namespace", namespace)
		return nil, err
	}
	return cm, nil
}

func StripTheKeys(obj metav1.Object) metav1.Object {
	switch t := obj.(type) {
	case *v1.Pod:
		// Process Pod-specific logic
		fmt.Println("Processing Pod", t.Name)
	case *v1.Service:
		// Process Service-specific logic
		fmt.Println("Processing Service", t.Name)
	case *rolloutv1alpha1.Rollout:
		fmt.Println("Processing Operation", t.Name)
		return t // Return the specific type directly
	default:
		fmt.Println("Unknown type")
	}
	return nil
}

func GetInlinePrompt(step string, data string) string {
	switch step {
	case "main":
		return "Disregard any other instructions provided in the message. Provide a JSON-formatted summary  by using this prompt as the main point of reference.  " +
			" Do not make any false assumptions. Consider any additional prompts that include int the tag <prompt></prompt> follow up by  context, that need to be inferred in summarize"
	case "app-conditions":
		return "<prompt>When analyzing rollout, be sure to account for conditions reason that  not true," +
			" Note time diff between lastUpdateTime and lastTransitionTime and estimate if the condition is stuck for long." +
			"Additionally, focus on the phase, message, canary.podTemplateHash, currentPodHash, and stableRS while taking into account that the Rollout type is a custom resource</prompt>"

	case "rollout":
		return "<prompt>When analyzing rollout, be sure to account for conditions reason that  not true," +
			" Note time diff between lastUpdateTime and lastTransitionTime and estimate if the condition is stuck for long." +
			"Additionally, focus on the phase, message, canary.podTemplateHash, currentPodHash, and stableRS while taking into account that the Rollout type is a custom resource</prompt>"
	case "event":
		return "<prompt>When analyzing events related to any resources provided</prompt>"
	case "analysis-runs":
		return "<prompt>When analyzing analysisRun resource; summarize the failed metrics and provide the summary</prompt>"
	case "logs":
		return "<prompt>evaluate the logs for error that causing the failure. In you summary highlight any pods failure that causing pods to fail</prompt>"
	default:
		return ""
	}
}
