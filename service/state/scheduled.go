package state

import (
	"fmt"
	v1 "github.com/Algatux/k8s-reconcyle-tests/api/v1"
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

func (o *OperationScheduled) Evolve(operation *v1.ScheduledOperation, r client.Writer) (ctrl.Result, error) {
	if !o.scheduler.MustBeExecuted(operation) {
		o.logger.Info(fmt.Sprintf(
			"OPERATION IS SCHEDULED, postponing execution to: %v",
			time.Unix(operation.Spec.NextExecutionTimestamp, 0),
		))
		return ctrl.Result{
			RequeueAfter: time.Duration(o.scheduler.SecondsToNextExecution(operation)) * time.Second,
		}, nil
	}

	o.logger.Info("OPERATION REACHED SCHEDULED EXECUTION TIME")
	operation.Spec.Status = v1.Ready

	return ctrl.Result{}, nil
}

func NewOperationScheduled(logger logr.Logger, scheduler *service.OperationScheduler) OperationStatus {
	return &OperationScheduled{
		baseOperationState: baseOperationState{logger: logger},
		scheduler:          scheduler,
	}
}
