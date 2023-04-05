package state

import (
	"context"
	v2 "github.com/Algatux/k8s-reconcyle-tests/api/v2"
	"github.com/Algatux/k8s-reconcyle-tests/service"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OperationFailure struct {
	baseOperationState
	scheduler *service.OperationScheduler
}

func (o *OperationFailure) Evolve(ctx context.Context, operation *v2.ScheduledOperation, r client.Writer) (ctrl.Result, error) {
	o.logger.Info("OPERATION FAILURE")
	if o.scheduler.MustReschedule(operation) {
		return ctrl.Result{}, o.scheduler.ScheduleOperation(operation)
	}

	return ctrl.Result{}, nil
}

func NewOperationFailure(logger logr.Logger, scheduler *service.OperationScheduler) OperationStatus {
	return &OperationFailure{
		baseOperationState: baseOperationState{logger: logger},
		scheduler:          scheduler,
	}
}
