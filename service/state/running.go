package state

import (
	"context"
	v2 "github.com/Algatux/k8s-reconcyle-tests/api/v2"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OperationRunning struct {
	baseOperationState
}

func (o *OperationRunning) Evolve(ctx context.Context, operation *v2.ScheduledOperation, r client.Writer) (ctrl.Result, error) {
	o.logger.Info("EXECUTING OPERATION")
	// DOING THINGS HERE TO KEEP THE TASK RUNNING

	// THEN CHANGE STATUS TO SUCCESS/FAILURE
	operation.Status.State = v2.Success
	operation.Status.Executions++

	return ctrl.Result{}, nil
}

func NewOperationRunning(logger logr.Logger) OperationStatus {
	return &OperationRunning{
		baseOperationState: baseOperationState{logger: logger},
	}
}
