package state

import (
	"context"
	v2 "github.com/Algatux/k8s-reconcyle-tests/api/v2"
	"github.com/Algatux/k8s-reconcyle-tests/service"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type OperationInit struct {
	baseOperationState
	ctx,
	scheduler *service.OperationScheduler
}

func (o *OperationInit) Evolve(ctx context.Context, operation *v2.ScheduledOperation, r client.Writer) (ctrl.Result, error) {
	o.logger.Info("OPERATION INITIALIZATION")

	controllerutil.AddFinalizer(operation, "ale/finalizer")
	operation.Labels["test"] = "ciao"
	r.Update(ctx, operation)

	if o.scheduler.IsScheduledOperation(operation) {
		return ctrl.Result{}, o.scheduler.InitScheduledOperation(operation)
	}

	operation.Status.State = v2.Ready

	return ctrl.Result{}, nil
}

func NewOperationInit(logger logr.Logger, scheduler *service.OperationScheduler) OperationStatus {
	return &OperationInit{
		baseOperationState: baseOperationState{logger: logger},
		scheduler:          scheduler,
	}
}
