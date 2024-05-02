package utils

import (
	"context"
	"fmt"
	v1alpha1 "github.com/argoproj-labs/argo-operations/api/v1alpha1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	// LabelKeyAppName is the label key to identify the authprovider
	LabelKeyAppName      = "app.kubernetes.io/name"
	LabelKeyAppNameValue = "argo-operations"
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

func GetConfigMap(ctx context.Context, k8sClient client.Client, obj metav1.Object) (*v1.ConfigMap, error) {
	logger := log.FromContext(ctx)

	var cm v1.ConfigMap
	objectKey := client.ObjectKey{
		Namespace: obj.GetNamespace(),
		Name:      obj.GetName(),
	}
	err := k8sClient.Get(ctx, objectKey, &cm)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("ConfigMap not found", "namespace", objectKey.Namespace, "name", objectKey.Name)
			return nil, err
		}

		logger.Error(err, "failed to get ConfigMap", "namespace", objectKey.Namespace, "name", objectKey.Name)
		return nil, err
	}

	return &cm, nil
}

func getAuthProvider(ctx context.Context, k8sClient client.Client, labels map[string]string) (*v1alpha1.AuthProvider, error) {
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

	return &authProviderList.Items[0], nil
}

func GetAIProviders(ctx context.Context, k8sClient client.Client, refs *[]v1alpha1.NamespacedObjectReference, obj metav1.Object) ([]*v1alpha1.AuthProvider, error) {
	logger := log.FromContext(ctx)

	var authProviders []*v1alpha1.AuthProvider
	for _, ref := range *refs {
		if authProvider, exists := authProviderMap[ref]; exists {
			// auth provider already fetched, skip processing
			authProviders = append(authProviders, authProvider)
			continue
		}

		// fetch auth provider
		authProvider, err := getAuthProvider(ctx, k8sClient, map[string]string{LabelKeyAppName: LabelKeyAppNameValue})
		if err != nil {
			logger.Error(err, "failed to get AuthProvider", "namespace", obj.GetNamespace(), "name", ref.Name)
			return nil, err
		}
		authProviderMap[ref] = authProvider

		authProviders = append(authProviders, authProvider)
	}

	return authProviders, nil
}
