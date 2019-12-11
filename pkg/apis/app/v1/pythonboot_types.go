package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PythonBoot is the Schema for the pythonboots API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:subresource:scale:specpath=.spec.replicas,statuspath=.status.HPAReplicas,selectorpath=.status.selector
// +kubebuilder:resource:path=pythonboots,shortName=python,scope=Namespaced
// +kubebuilder:printcolumn:name="Desired",type="integer",JSONPath=".spec.replicas",description="Number of desired pods"
// +kubebuilder:printcolumn:name="ReadyReplicas",type="integer",JSONPath=".status.readyReplicas",description="Number of ready pods"
// +kubebuilder:printcolumn:name="CurrentReplicas",type="integer",JSONPath=".status.currentReplicas",description="Number of current pods"
// +kubebuilder:printcolumn:name="Services",type="string",JSONPath=".status.services",description="The service's name of the boot"
// +kubebuilder:printcolumn:name="Workload",type="string",JSONPath=".status.workload",description="The wordload type for the boot"
// +kubebuilder:printcolumn:name="Version",type="string",JSONPath=".spec.version",description="The Version of Boot"
type PythonBoot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BootSpec   `json:"spec,omitempty"`
	Status BootStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PythonBootList contains a list of PythonBoot
type PythonBootList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PythonBoot `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PythonBoot{}, &PythonBootList{})
}
