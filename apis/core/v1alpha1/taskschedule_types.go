/*
Copyright 2022 The Crossplane Authors.

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

package v1alpha1

import (
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// TaskScheduleParameters are the configurable fields of a TaskSchedule.
type TaskScheduleParameters struct {

	// Name of the task schedule (immutable)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:MaxLength=52
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Name is immutable"
	ScheduleName string `json:"scheduleName"`

	// TaskDefinition Name that will be scheduled (immutable)
	// At least one of taskDefinitionName, taskDefinitionNameRef or taskDefinitionNameSelector is required.
	// +kubebuilder:validation:Optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="TaskDefinitionName is immutable"
	// +crossplane:generate:reference:type=github.com/denniskniep/provider-springclouddataflow/apis/core/v1alpha1.TaskDefinition
	TaskDefinitionName *string `json:"taskDefinitionName,omitempty"`

	// TaskDefinition reference to retrieve the TaskDefinition Name, that will be scheduled
	// At least one of taskDefinitionName, taskDefinitionNameRef or taskDefinitionNameSelector is required.
	// +optional
	TaskDefinitionNameRef *xpv1.Reference `json:"taskDefinitionNameRef,omitempty"`

	// TaskDefinitionNameSelector selects a reference to a TaskDefinition and retrieves its name
	// At least one of temporalNamespaceName, temporalNamespaceNameRef or temporalNamespaceNameSelector is required.
	// +optional
	TaskDefinitionNameSelector *xpv1.Selector `json:"taskDefinitionNameSelector,omitempty"`

	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Platform is immutable"
	// +kubebuilder:validation:Required
	CronExpression string `json:"cronExpression,omitempty"`

	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Platform is immutable"
	// +kubebuilder:default=default
	Platform string `json:"platform,omitempty"`

	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Arguments is immutable"
	Arguments *string `json:"arguments,omitempty"`

	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Properties is immutable"
	Properties *string `json:"properties,omitempty"`
}

// TaskScheduleObservation are the observable fields of a TaskSchedule.
type TaskScheduleObservation struct {
	ScheduleName string `json:"scheduleName"`

	TaskDefinitionName *string `json:"taskDefinitionName,omitempty"`
}

// A TaskScheduleSpec defines the desired state of a TaskSchedule.
type TaskScheduleSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       TaskScheduleParameters `json:"forProvider"`
}

// A TaskScheduleStatus represents the observed state of a TaskSchedule.
type TaskScheduleStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          TaskScheduleObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A TaskSchedule is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,springclouddataflow}
type TaskSchedule struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TaskScheduleSpec   `json:"spec"`
	Status TaskScheduleStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TaskScheduleList contains a list of TaskSchedule
type TaskScheduleList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TaskSchedule `json:"items"`
}

// TaskSchedule type metadata.
var (
	TaskScheduleKind             = reflect.TypeOf(TaskSchedule{}).Name()
	TaskScheduleGroupKind        = schema.GroupKind{Group: Group, Kind: TaskScheduleKind}.String()
	TaskScheduleKindAPIVersion   = TaskScheduleKind + "." + SchemeGroupVersion.String()
	TaskScheduleGroupVersionKind = SchemeGroupVersion.WithKind(TaskScheduleKind)
)

func init() {
	SchemeBuilder.Register(&TaskSchedule{}, &TaskScheduleList{})
}
