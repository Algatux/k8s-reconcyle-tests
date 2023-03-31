package state

import (
	v1 "github.com/Algatux/k8s-reconcyle-tests/api/v1"
	"github.com/Algatux/k8s-reconcyle-tests/service"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OperationInit struct {
	baseOperationState
	scheduler *service.OperationScheduler
}

func (o *OperationInit) Evolve(operation *v1.ScheduledOperation, r client.Writer) (ctrl.Result, error) {
	o.logger.Info("OPERATION INITIALIZATION")
	if o.scheduler.IsScheduledOperation(operation) {
		return ctrl.Result{}, o.scheduler.InitScheduledOperation(operation)
	}

	operation.Spec.Status = v1.Ready

	return ctrl.Result{}, nil
}

func NewOperationInit(logger logr.Logger, scheduler *service.OperationScheduler) OperationStatus {
	return &OperationInit{
		baseOperationState: baseOperationState{logger: logger},
		scheduler:          scheduler,
	}
}
