package v1alpha1

import (
	"github.com/argoproj/argo-rollouts/pkg/apis/rollouts"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: rollouts.Group, Version: "v1alpha1"}

var (
	// GroupVersionResource for all rollout types
	RolloutGVR = SchemeGroupVersion.WithResource("rollouts")
)

const (
	ResourcesFinalizerName string = "resources-finalizer.argocd.argoproj.io"

	// LabelKeyAppName is the label key to identify the authprovider
	LabelKeyAppName      = "app.kubernetes.io/name"
	LabelKeyAppNameValue = "argo-support"
	FinalizerName        = "support.argoproj.extensions.io/finalizer"
)

type ArgoSupportPhase string

// Possible ArgoSupportPhase values
const (
	ArgoSupportPhaseRunning   ArgoSupportPhase = "Running"
	ArgoSupportPhaseCompleted ArgoSupportPhase = "Completed"
	ArgoSupportPhaseFailed    ArgoSupportPhase = "Failed"
)

type Auth struct {
	BaseURL          string `json:"baseUrl,omitempty"`
	AppID            string `json:"appId,omitempty"`
	IdentityEndpoint string `json:"identityEndpoint,omitempty"`
	IdentityJobID    string `json:"identityJobID,omitempty"`
	APIVersion       string `json:"apiVersion,omitempty"`
}

type Workflow struct {

	// +kubebuilder:validation:Required
	Name string `json:"name"`
	// +kubebuilder:validation:Required
	Initiate bool `json:"initiate"`
	// +kubebuilder:validation:Required
	Ref []NamespacedObjectReference `json:"autProviderRef"`
	// +kubebuilder:validation:Optional
	ConfigMapRef ConfigMapRef `json:"configMapRef"`
}

type NamespacedObjectReference struct {
	// +kubebuilder:validation:Required
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}

type ConfigMapRef struct {
	// Name of the ConfigMap
	Name string `json:"name"`
}
