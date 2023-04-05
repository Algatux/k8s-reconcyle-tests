package state

import (
	"errors"
	v2 "github.com/Algatux/k8s-reconcyle-tests/api/v2"
	"github.com/Algatux/k8s-reconcyle-tests/service"
	"github.com/go-logr/logr"
)

type OperationsStateFactory struct {
	logger    logr.Logger
	scheduler *service.OperationScheduler
}

func (osf *OperationsStateFactory) GetStateByOperation(operation *v2.ScheduledOperation) (OperationStatus, error) {
	if operation.Status.State == "" {
		operation.Status.State = v2.Init
	}

	switch operation.Status.State {
	case v2.Init:
		return NewOperationInit(osf.logger, osf.scheduler), nil
	case v2.Scheduled:
		return NewOperationScheduled(osf.logger, osf.scheduler), nil
	case v2.Ready:
		return NewOperationReady(osf.logger), nil
	case v2.Running:
		return NewOperationRunning(osf.logger), nil
	case v2.Success:
		return NewOperationSuccess(osf.logger, osf.scheduler), nil
	case v2.Failure:
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
