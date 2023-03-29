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
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"time"

	operationsv1 "github.com/Algatux/k8s-reconcyle-tests/api/v1"
)

// ScheduledOperationReconciler reconciles a ScheduledOperation object
type ScheduledOperationReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	logger logr.Logger
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
	r.logger = log.FromContext(ctx)

	var operation operationsv1.ScheduledOperation
	if err := r.Get(ctx, req.NamespacedName, &operation); err != nil {
		r.logger.Error(err, "unable to fetch ScheduledOperation")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if operation.Spec.Status == operationsv1.Init {
		r.logger.Info("OPERATION INITIALIZATION")
		operation.Spec.Status = operationsv1.Scheduled

		return ctrl.Result{}, r.updateOperation(ctx, &operation)
	}

	missingSeconds := operation.CreationTimestamp.Unix() + 60 - time.Now().Unix()
	if operation.Spec.Status == operationsv1.Scheduled && missingSeconds > 0 {
		r.logger.Info(fmt.Sprintf("OPERATION IS SCHEDULED, time to execution: %d", missingSeconds))

		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	r.logger.Info("EXECUTING OPERATION ")

	return ctrl.Result{}, nil
}

func (r *ScheduledOperationReconciler) updateOperation(ctx context.Context, operation *operationsv1.ScheduledOperation) error {
	err := r.Update(ctx, operation)
	if err != nil {
		r.logger.Error(err, "update failed")
		return err
	}

	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *ScheduledOperationReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operationsv1.ScheduledOperation{}).
		Complete(r)
}
