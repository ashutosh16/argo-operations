/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"github.com/argoproj-labs/argo-operations/internal/wf_operation/ai_operations"
	"github.com/argoproj-labs/argo-operations/internal/wf_operation/ai_operations/genai"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	argosupportv1alpha1 "github.com/argoproj-labs/argo-operations/api/v1alpha1"
)

// ArgoAISupportReconciler reconciles a ArgoAISupport object
type ArgoAISupportReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	DynamicClient *dynamic.DynamicClient
}

//+kubebuilder:rbac:groups=argosupport.argoproj.extensions.io,resources=argoaisupports,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=argosupport.argoproj.extensions.io,resources=argoaisupports/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=argosupport.argoproj.extensions.io,resources=argoaisupports/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ArgoAISupport object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.2/pkg/reconcile
func (r *ArgoAISupportReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var aiSupport argosupportv1alpha1.ArgoAISupport
	err := r.Get(ctx, req.NamespacedName, &aiSupport, &client.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Info("argo AI support operation not found", "namespace", req.Namespace, "name", req.Name)
			return ctrl.Result{}, nil
		}

		logger.Info("argo AI support operation not found", "namespace", req.Namespace, "name", req.Name)
		return ctrl.Result{}, err
	}

	wf, err := r.getWfExecutor(ctx, &aiSupport.Spec.Workflows[0], &aiSupport)
	//logger.Error(fmt.Errorf("argo AI support operation not found", "namespace", req.Namespace, "name", req.Name))

	_, err = wf.Process(ctx, &aiSupport)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.handleFinalizer(ctx, &aiSupport)
	if err != nil {
		return ctrl.Result{}, err
	}

	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ArgoAISupportReconciler) getWfExecutor(ctx context.Context, wf *argosupportv1alpha1.Workflow, obj metav1.Object) (ai_operations.Executor, error) {

	switch {
	case wf.Name == "gen-ai":
		genaiops, err := genai.NewGenAIOperations(ctx, r.Client, r.DynamicClient, wf, obj)
		if err != nil {
			return nil, err
		}
		return genaiops, nil
	default:
		return nil, nil
	}
}

func (r *ArgoAISupportReconciler) handleFinalizer(ctx context.Context, ops *argosupportv1alpha1.ArgoAISupport) error {
	// name of our custom finalizer
	finalizerName := "support.argoproj.extensions.io/finalizer"

	// examine DeletionTimestamp to determine if object is under deletion
	if ops.ObjectMeta.DeletionTimestamp.IsZero() {
		// The object is not being deleted, so if it does not have our finalizer,
		// then lets add the finalizer and update the object. This is equivalent
		// to registering our finalizer.
		if !controllerutil.ContainsFinalizer(ops, finalizerName) {
			controllerutil.AddFinalizer(ops, finalizerName)
			if err := r.Update(ctx, ops); err != nil {
				return err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(ops, finalizerName) {
			// our finalizer is present, so lets handle any external dependency
			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(ops, finalizerName)
			if err := r.Update(ctx, ops); err != nil {
				return err
			}
		}
	}

	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ArgoAISupportReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&argosupportv1alpha1.ArgoAISupport{}).
		Complete(r)
}
