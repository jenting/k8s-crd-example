package v1

import (
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Health describes a Health resource
type Health struct {
	// TypeMeta is the metadata for the resource, like kind and apiversion
	metaV1.TypeMeta `json:",inline"`
	// ObjectMeta contains the metadata for the particular object
	metaV1.ObjectMeta `json:"metadata,omitempty"`

	// Spec is the custom resource spec
	Spec HealthSpec `json:"spec"`
}

// HealthSpec is the spec for a Health resource
type HealthSpec struct {
	Action string `json:"action"`
	Switch bool   `json:"switch"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// HealthList is a list of Health resources
type HealthList struct {
	metaV1.TypeMeta `json:",inline"`
	metaV1.ListMeta `json:"metadata"`

	Items []Health `json:"items"`
}
