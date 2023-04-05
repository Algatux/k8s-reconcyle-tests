package state

import (
	"context"
	v2 "github.com/Algatux/k8s-reconcyle-tests/api/v2"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OperationReady struct {
	baseOperationState
}

func (o *OperationReady) Evolve(ctx context.Context, operation *v2.ScheduledOperation, r client.Writer) (ctrl.Result, error) {
	o.logger.Info("OPERATION READY")
	operation.Status.State = v2.Running
	// DOING THINGS HERE TO START THE TASKS

	return ctrl.Result{}, nil
}

func NewOperationReady(logger logr.Logger) OperationStatus {
	return &OperationReady{
		baseOperationState: baseOperationState{logger: logger},
	}
}
