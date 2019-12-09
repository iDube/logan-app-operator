package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BootRevision is the Schema for the bootrevisions API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=bootrevisions,scope=Namespaced
// +kubebuilder:printcolumn:name="Desired",type="integer",JSONPath=".spec.replicas",description="Number of desired pods"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version",description="The Version of Boot"
type BootRevision struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec contains the desired behavior of the Boot
	Spec BootSpec `json:"spec,omitempty"`
	// status contains the last observed state of the BootStatus
	Status BootStatus `json:"status,omitempty"`

	BootType string `json:"bootType"`
	AppKey   string `json:"appKey"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BootRevisionList contains a list of BootRevision
type BootRevisionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []BootRevision `json:"items"`
}

func init() {
	SchemeBuilder.Register(&BootRevision{}, &BootRevisionList{})
}
