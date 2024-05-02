package ai_operations

import (
	"context"
	"github.com/argoproj-labs/argo-operations/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Executor interface {
	Process(ctx context.Context, obj metav1.Object) (*v1alpha1.ArgoAISupport, error)
}
