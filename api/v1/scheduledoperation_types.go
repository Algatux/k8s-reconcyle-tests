/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type OperationState string

const (
	Init      OperationState = "INIT"
	Scheduled OperationState = "SCHEDULED"
	Running   OperationState = "RUNNING"
	Ready     OperationState = "READY"
	Success   OperationState = "SUCCESS"
	Failure   OperationState = "FAILURE"
)

// ScheduledOperationSpec defines the desired state of ScheduledOperation
type ScheduledOperationSpec struct {
	// The schedule in Cron format, see https://en.wikipedia.org/wiki/Cron.
	// +kubebuilder:validation:MinLength=0
	// +kubebuilder:default=""
	// +optional
	Schedule string `json:"schedule"`
	// Number of times the operation must be executed on schedule
	// +kubebuilder:default=-1
	// +kubebuilder:validation:Minimum=-1
	// +optional
	DesiredExecutions int `json:"desiredExecutions"`
}

// ScheduledOperationStatus defines the observed state of ScheduledOperation
type ScheduledOperationStatus struct {
	// State of the operation
	// +kubebuilder:default=INIT
	// +kubebuilder:validation:Enum=INIT;SCHEDULED;READY;RUNNING;SUCCESS;FAILURE
	// +optional
	State OperationState `json:"state"`
	// Number of times the operation has been executed on schedule
	// +kubebuilder:default=0
	// +kubebuilder:validation:Minimum=0
	// +optional
	Executions int `json:"executions"`
	// +kubebuilder:default=0
	// +optional
	NextExecutionTimestamp int64 `json:"nextExecutionTimestamp"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ScheduledOperation is the Schema for the scheduledoperations API
type ScheduledOperation struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ScheduledOperationSpec   `json:"spec,omitempty"`
	Status ScheduledOperationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ScheduledOperationList contains a list of ScheduledOperation
type ScheduledOperationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ScheduledOperation `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ScheduledOperation{}, &ScheduledOperationList{})
}
