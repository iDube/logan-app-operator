package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeJSBoot is the Schema for the nodejsboots API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=nodejsboots,shortName=nodejs,scope=Namespaced
// +kubebuilder:printcolumn:name="Desired",type="integer",JSONPath=".spec.replicas",description="Number of desired pods"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version",description="The Version of Boot"
type NodeJSBoot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BootSpec   `json:"spec,omitempty"`
	Status BootStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeJSBootList contains a list of NodeJSBoot
type NodeJSBootList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodeJSBoot `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodeJSBoot{}, &NodeJSBootList{})
}
