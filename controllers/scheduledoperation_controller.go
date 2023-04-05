/*
Copyright 2023.

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

package controllers

import (
	"context"
	"fmt"
	operationsv2 "github.com/Algatux/k8s-reconcyle-tests/api/v2"
	"github.com/Algatux/k8s-reconcyle-tests/service/state"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ScheduledOperationReconciler reconciles a ScheduledOperation object
type ScheduledOperationReconciler struct {
	client.Client
	Scheme       *runtime.Scheme
	Logger       logr.Logger
	StateFactory *state.OperationsStateFactory
}

//+kubebuilder:rbac:groups=operations.algatux.dev,resources=scheduledoperations,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operations.algatux.dev,resources=scheduledoperations/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operations.algatux.dev,resources=scheduledoperations/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ScheduledOperation object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *ScheduledOperationReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var err error
	var operation operationsv2.ScheduledOperation
	if err = r.Get(ctx, req.NamespacedName, &operation); err != nil {
		r.Logger.Error(err, "unable to fetch ScheduledOperation")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	r.Logger.Info(fmt.Sprintf("|| >> Operation Reconcile cycle: %s", operation.Name))
	r.Logger.Info(fmt.Sprintf("Status: %s", operation.Status.State))

	operationState, err := r.StateFactory.GetStateByOperation(&operation)
	if err != nil {
		r.Logger.Error(err, "Error occurred during operation state evaluation")
		return ctrl.Result{}, err
	}

	result, err := operationState.Evolve(ctx, &operation, r)
	if err != nil {
		r.Logger.Error(err, "Error evolving operation state")
		return ctrl.Result{}, err
	}

	return result, r.Status().Update(ctx, &operation)
}

// SetupWithManager sets up the controller with the Manager.
func (r *ScheduledOperationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operationsv2.ScheduledOperation{}).
		Complete(r)
}
