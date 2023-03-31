package state

import (
	v1 "github.com/Algatux/k8s-reconcyle-tests/api/v1"
	"github.com/go-logr/logr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type OperationStatus interface {
	Evolve(operation *v1.ScheduledOperation, r client.Writer) (ctrl.Result, error)
}

type baseOperationState struct {
	logger logr.Logger
}
