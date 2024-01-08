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

// TaskDefinitionParameters are the configurable fields of a TaskDefinition.
type TaskDefinitionParameters struct {

	// Name of the task definition (immutable)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Name is immutable"
	Name string `json:"name"`

	// Description of the task definition (immutable)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Description is immutable"
	Description string `json:"description"`

	// The definition for the task, using Data Flow DSL (immutable)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Definition is immutable"
	Definition string `json:"definition"`
}

// TaskDefinitionObservation are the observable fields of a TaskDefinition.
type TaskDefinitionObservation struct {
	Name                string `json:"name"`
	Description         string `json:"description"`
	Definition          string `json:"definition"`
	Composed            bool   `json:"composed"`
	ComposedTaskElement bool   `json:"composedTaskElement"`
	Status              string `json:"status"`
}

// A TaskDefinitionSpec defines the desired state of a TaskDefinition.
type TaskDefinitionSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       TaskDefinitionParameters `json:"forProvider"`
}

// A TaskDefinitionStatus represents the observed state of a TaskDefinition.
type TaskDefinitionStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          TaskDefinitionObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A TaskDefinition is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,springclouddataflow}
type TaskDefinition struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TaskDefinitionSpec   `json:"spec"`
	Status TaskDefinitionStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// TaskDefinitionList contains a list of TaskDefinition
type TaskDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []TaskDefinition `json:"items"`
}

// TaskDefinition type metadata.
var (
	TaskDefinitionKind             = reflect.TypeOf(TaskDefinition{}).Name()
	TaskDefinitionGroupKind        = schema.GroupKind{Group: Group, Kind: TaskDefinitionKind}.String()
	TaskDefinitionKindAPIVersion   = TaskDefinitionKind + "." + SchemeGroupVersion.String()
	TaskDefinitionGroupVersionKind = SchemeGroupVersion.WithKind(TaskDefinitionKind)
)

func init() {
	SchemeBuilder.Register(&TaskDefinition{}, &TaskDefinitionList{})
}
