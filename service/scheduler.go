package service

import (
	"errors"
	"fmt"
	v2 "github.com/Algatux/k8s-reconcyle-tests/api/v2"
	"github.com/go-logr/logr"
	"github.com/robfig/cron/v3"
	"time"
)

const (
	AlwaysRepeat = -1
)

type OperationScheduler struct {
	logger logr.Logger
	parser cron.Parser
}

func NewScheduler(logger logr.Logger) *OperationScheduler {
	return &OperationScheduler{
		logger: logger,
		parser: cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow),
	}
}

func (s *OperationScheduler) IsScheduledOperation(operation *v2.ScheduledOperation) bool {
	return len(operation.Spec.Schedule) > 0
}

func (s *OperationScheduler) InitScheduledOperation(operation *v2.ScheduledOperation) error {
	if operation.Status.NextExecutionTimestamp != 0 {
		return errors.New("operation already initialized")
	}
	return s.ScheduleOperation(operation)
}

func (s *OperationScheduler) ScheduleOperation(operation *v2.ScheduledOperation) error {
	s.logger.Info("SCHEDULING OPERATION")
	operation.Status.State = v2.Scheduled
	s.logger.Info(fmt.Sprintf("OPERATION SCHEDULE : %v", operation.Spec.Schedule))
	nextExecution, err := s.getNextExecution(operation)
	if err != nil {
		s.logger.Error(err, "Error parsing operation schedule")
		return err
	}
	s.logger.Info(fmt.Sprintf(
		"OPERATION IS SCHEDULED RUN %d of %d, next execution: %v",
		operation.Status.Executions+1,
		operation.Spec.DesiredExecutions,
		nextExecution,
	))
	operation.Status.NextExecutionTimestamp = nextExecution.Unix()

	return nil
}

func (s *OperationScheduler) SecondsToNextExecution(operation *v2.ScheduledOperation) int64 {
	return operation.Status.NextExecutionTimestamp - time.Now().Unix()
}

func (s *OperationScheduler) MustBeExecuted(operation *v2.ScheduledOperation) bool {
	return s.SecondsToNextExecution(operation) <= 0
}

func (s *OperationScheduler) MustReschedule(operation *v2.ScheduledOperation) bool {
	if !s.IsScheduledOperation(operation) {
		return false
	}

	return operation.Spec.DesiredExecutions == AlwaysRepeat || operation.Status.Executions < operation.Spec.DesiredExecutions
}

func (s *OperationScheduler) getNextExecution(operation *v2.ScheduledOperation) (*time.Time, error) {
	schedule, err := s.parser.Parse(operation.Spec.Schedule)
	if err != nil {
		return nil, err
	}
	nextExecution := schedule.Next(time.Now())
	return &nextExecution, err
}
