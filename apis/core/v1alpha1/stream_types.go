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

// StreamParameters are the configurable fields of a Stream.
type StreamParameters struct {
	// Name of the stream (immutable)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Name is immutable"
	Name string `json:"name"`

	// Description of the stream (immutable)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Description is immutable"
	Description string `json:"description"`

	// The definition for the stream, using Data Flow DSL (immutable)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Definition is immutable"
	Definition string `json:"definition"`

	// If true, the stream is deployed upon creation (immutable)
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Deploy is immutable"
	Deploy bool `json:"deploy"`
}

// StreamObservation are the observable fields of a Stream.
type StreamObservation struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	Definition        string `json:"definition"`
	Status            string `json:"status"`
	StatusDescription string `json:"statusDescription"`
}

// A StreamSpec defines the desired state of a Stream.
type StreamSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       StreamParameters `json:"forProvider"`
}

// A StreamStatus represents the observed state of a Stream.
type StreamStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          StreamObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Stream is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,springclouddataflow}
type Stream struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StreamSpec   `json:"spec"`
	Status StreamStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// StreamList contains a list of Stream
type StreamList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Stream `json:"items"`
}

// Stream type metadata.
var (
	StreamKind             = reflect.TypeOf(Stream{}).Name()
	StreamGroupKind        = schema.GroupKind{Group: Group, Kind: StreamKind}.String()
	StreamKindAPIVersion   = StreamKind + "." + SchemeGroupVersion.String()
	StreamGroupVersionKind = SchemeGroupVersion.WithKind(StreamKind)
)

func init() {
	SchemeBuilder.Register(&Stream{}, &StreamList{})
}
