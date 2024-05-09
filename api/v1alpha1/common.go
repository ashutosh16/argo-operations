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
}

type NamespacedObjectReference struct {
	// +kubebuilder:validation:Required
	Name      string `json:"name"`
	Namespace string `json:"namespace,omitempty"`
}
