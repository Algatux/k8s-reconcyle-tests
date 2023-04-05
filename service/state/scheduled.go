package state

import (
	"context"
	"fmt"
	v2 "github.com/Algatux/k8s-reconcyle-tests/api/v2"
	"github.com/Algatux/k8s-reconcyle-tests/service"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type OperationScheduled struct {
	baseOperationState
	scheduler *service.OperationScheduler
}

func (o *OperationScheduled) Evolve(ctx context.Context, operation *v2.ScheduledOperation, r client.Writer) (ctrl.Result, error) {
	if !o.scheduler.MustBeExecuted(operation) {
		o.logger.Info(fmt.Sprintf(
			"OPERATION IS SCHEDULED, postponing execution to: %v",
			time.Unix(operation.Status.NextExecutionTimestamp, 0),
		))
		return ctrl.Result{
			RequeueAfter: time.Duration(o.scheduler.SecondsToNextExecution(operation)) * time.Second,
		}, nil
	}

	o.logger.Info("OPERATION REACHED SCHEDULED EXECUTION TIME")
	operation.Status.State = v2.Ready

	return ctrl.Result{}, nil
}

func NewOperationScheduled(logger logr.Logger, scheduler *service.OperationScheduler) OperationStatus {
	return &OperationScheduled{
		baseOperationState: baseOperationState{logger: logger},
		scheduler:          scheduler,
	}
}
