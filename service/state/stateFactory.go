package state

import (
	"errors"
	v1 "github.com/Algatux/k8s-reconcyle-tests/api/v1"
	"github.com/Algatux/k8s-reconcyle-tests/service"
	"github.com/go-logr/logr"
)

type OperationsStateFactory struct {
	logger    logr.Logger
	scheduler *service.OperationScheduler
}

func (osf *OperationsStateFactory) GetStateByOperation(operation *v1.ScheduledOperation) (OperationStatus, error) {
	switch operation.Spec.Status {
	case v1.Init:
		return NewOperationInit(osf.logger, osf.scheduler), nil
	case v1.Scheduled:
		return NewOperationScheduled(osf.logger, osf.scheduler), nil
	case v1.Ready:
		return NewOperationReady(osf.logger), nil
	case v1.Running:
		return NewOperationRunning(osf.logger), nil
	case v1.Success:
		return NewOperationSuccess(osf.logger, osf.scheduler), nil
	case v1.Failure:
		return NewOperationFailure(osf.logger, osf.scheduler), nil
	}

	return nil, errors.New("no valid state received")
}

func NewOperationStateFactory(logger logr.Logger, scheduler *service.OperationScheduler) *OperationsStateFactory {
	return &OperationsStateFactory{
		logger:    logger,
		scheduler: scheduler,
	}
}
