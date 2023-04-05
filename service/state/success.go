package state

import (
	"context"
	v2 "github.com/Algatux/k8s-reconcyle-tests/api/v2"
	"github.com/Algatux/k8s-reconcyle-tests/service"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OperationSuccess struct {
	baseOperationState
	scheduler *service.OperationScheduler
}

func (o *OperationSuccess) Evolve(ctx context.Context, operation *v2.ScheduledOperation, r client.Writer) (ctrl.Result, error) {
	o.logger.Info("OPERATION SUCCESS")
	if o.scheduler.MustReschedule(operation) {
		return ctrl.Result{}, o.scheduler.ScheduleOperation(operation)
	}

	return ctrl.Result{}, nil
}

func NewOperationSuccess(logger logr.Logger, scheduler *service.OperationScheduler) OperationStatus {
	return &OperationSuccess{
		baseOperationState: baseOperationState{logger: logger},
		scheduler:          scheduler,
	}
}
