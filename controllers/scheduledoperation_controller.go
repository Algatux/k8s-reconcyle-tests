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
	"github.com/Algatux/k8s-reconcyle-tests/service"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	operationsv1 "github.com/Algatux/k8s-reconcyle-tests/api/v1"
)

// ScheduledOperationReconciler reconciles a ScheduledOperation object
type ScheduledOperationReconciler struct {
	client.Client
	Scheme    *runtime.Scheme
	Logger    logr.Logger
	Scheduler service.OperationScheduler
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
	var operation operationsv1.ScheduledOperation
	if err = r.Get(ctx, req.NamespacedName, &operation); err != nil {
		r.Logger.Error(err, "unable to fetch ScheduledOperation")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	r.Logger.Info(fmt.Sprintf("Operation Reconcile: %s", operation.Name))

	if operation.Spec.Status == operationsv1.Init {
		err := r.initOperation(&operation)
		if err != nil {
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, r.Update(ctx, &operation)
	}

	if operation.Spec.Status == operationsv1.Scheduled {
		if !r.Scheduler.MustBeExecuted(&operation) {
			r.Logger.Info(fmt.Sprintf(
				"OPERATION IS SCHEDULED, postponing execution to: %v",
				time.Unix(operation.Spec.NextExecution, 0),
			))
			return ctrl.Result{
				RequeueAfter: time.Duration(r.Scheduler.SecondsToNextExecution(&operation)) * time.Second,
			}, nil
		}

		r.makeOperationReady(&operation)
		return ctrl.Result{}, r.updateOperation(ctx, &operation)
	}

	if operation.Spec.Status == operationsv1.Ready {
		operation.Spec.Status = operationsv1.Running
		return ctrl.Result{}, r.updateOperation(ctx, &operation)
	}

	if operation.Spec.Status == operationsv1.Running {
		r.Logger.Info("EXECUTING OPERATION ")
		operation.Spec.Status = operationsv1.Success
		operation.Spec.Executions++
	}

	if operation.Spec.Status == operationsv1.Success {
		r.Logger.Info("OPERATION SUCCESS")
		if r.Scheduler.MustReschedule(&operation) {
			err := r.Scheduler.ScheduleOperation(&operation)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	return ctrl.Result{}, r.updateOperation(ctx, &operation)
}

func (r *ScheduledOperationReconciler) makeOperationReady(operation *operationsv1.ScheduledOperation) {
	operation.Spec.Status = operationsv1.Ready
	r.Logger.Info("OPERATION READY")
}

func (r *ScheduledOperationReconciler) initOperation(operation *operationsv1.ScheduledOperation) error {
	r.Logger.Info("OPERATION INITIALIZATION")
	if r.Scheduler.IsScheduledOperation(operation) {
		return r.Scheduler.InitScheduledOperation(operation)
	}

	r.makeOperationReady(operation)

	return nil
}

func (r *ScheduledOperationReconciler) updateOperation(ctx context.Context, operation *operationsv1.ScheduledOperation) error {
	err := r.Update(ctx, operation)
	if err != nil {
		r.Logger.Error(err, "update failed")
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
