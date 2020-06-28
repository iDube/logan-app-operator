package v1

import (
	autoscaling "k8s.io/api/autoscaling/v2beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Boot is the common Schema for the all boot types API
type Boot struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec contains the desired behavior of the Boot
	Spec BootSpec `json:"spec,omitempty"`
	// status contains the last observed state of the BootStatus
	Status BootStatus `json:"status,omitempty"`

	BootType string `json:"bootType"`
	AppKey   string `json:"appKey"`
}

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BootSpec defines the desired state of Boot for specified types, as JavaBoot/PhpBoot/PythonBoot/NodeJSBoot
type BootSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// Image is the app container' image. Image must not have a tag version.
	Image string `json:"image"`
	// Version is the app container's image version.
	Version string `json:"version"`
	// Replicas is the number of desired replicas.
	// This is a pointer to distinguish between explicit zero and unspecified.
	// Defaults to 1.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=100
	Replicas *int32 `json:"replicas,omitempty"`
	// Env is list of environment variables to set in the app container.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Env []corev1.EnvVar `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
	// Port that are exposed by the app container
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port,omitempty"`
	// Reserved, not used. for latter use
	SubDomain string `json:"subDomain,omitempty"`
	// Health is check path for the app container.
	// +kubebuilder:validation:MinLength=0
	// +kubebuilder:validation:MaxLength=2048
	Health *string `json:"health,omitempty"`
	// Readiness is a readiness check path for the app container.
	// +kubebuilder:validation:MinLength=0
	// +kubebuilder:validation:MaxLength=2048
	Readiness *string `json:"readiness,omitempty"`
	// Prometheus will scrape metrics from the service, default is `true`
	// +kubebuilder:validation:Enum="true";"false";""
	Prometheus string `json:"prometheus,omitempty"`
	// Resources is the compute resource requirements for the app container
	Resources corev1.ResourceRequirements `json:"resources,omitempty"`
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// Command is command for boot's container. If empty, will use image's ENTRYPOINT, specified here if needed override.
	Command []string `json:"command,omitempty"`
	// SessionAffinity is SessionAffinity for boot's created service. If empty, will not set
	// +kubebuilder:validation:Enum=ClientIP;None
	SessionAffinity string `json:"sessionAffinity,omitempty"`
	// NodePort will expose the service on each node’s IP at a random port, default is ``
	// +kubebuilder:validation:Enum=true;false
	NodePort string `json:"nodePort,omitempty"`
	// pvc is list of PersistentVolumeClaim to set in the app container.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Pvc []PersistentVolumeClaimMount `json:"pvc,omitempty" patchStrategy:"merge" patchMergeKey:"name"`
	// Priority will set the priorityClassName for the boot's workloads, default is ``
	Priority string `json:"priority,omitempty"`
	// Workload will set the wordload type for the boot,can be `Deployment` or `StatefulSet`. default is `Deployment`
	// +kubebuilder:validation:Enum=Deployment;StatefulSet
	Workload Workload `json:"workload,omitempty"`
	// Hpa is the configuration for a horizontal pod
	// autoscaler, which automatically manages the replica count of any resource
	// implementing the scale subresource based on the metrics specified.
	// +optional
	Hpa *Hpa `json:"hpa,omitempty"`
}

// Workload defines the wordload type for the boot
type Workload string

const (
	// Deployment defines the Deployment wordload type
	Deployment Workload = "Deployment"
	// StatefulSet defines the StatefulSet wordload type
	StatefulSet Workload = "StatefulSet"
)

// BootStatus defines the observed state of Boot for specified types, as JavaBoot/PhpBoot/PythonBoot/NodeJSBoot
type BootStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html

	// Services is the service's name of the boot, include app and sidecar
	// +optional
	Services string `json:"services,omitempty"`
	// Workload is the wordload type for the boot,can be `Deployment` or `StatefulSet`
	// +optional
	// +kubebuilder:validation:Enum=Deployment;StatefulSet
	Workload Workload `json:"workload,omitempty"`
	// HPAReplicas the number of non-terminated replicas that are receiving active traffic
	// +optional
	HPAReplicas int32 `json:"HPAReplicas,omitempty"`
	// Selector that identifies the pods that are receiving active traffic
	// +optional
	Selector string `json:"selector,omitempty"`
	// Replicas is the number of desired replicas.
	// +optional
	Replicas int32 `json:"replicas,omitempty"`
	// CurrentReplicas is the number of current replicas.
	// +optional
	CurrentReplicas int32 `json:"currentReplicas,omitempty"`
	// ReadyReplicas is the number of ready replicas.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty"`
	// Revision is the revision ID of the boot
	// +optional
	Revision string `json:"revision,omitempty"`
}

// PersistentVolumeClaimMount defines the Boot match a PersistentVolumeClaim
type PersistentVolumeClaimMount struct {
	// This must match the Name of a PersistentVolumeClaim.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Mounted read-only if true, read-write otherwise (false or unspecified).
	// Defaults to false.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,2,opt,name=readOnly"`
	// Path within the container at which the volume should be mounted.  Must
	// not contain ':'.
	// +kubebuilder:validation:MinLength=1
	MountPath string `json:"mountPath" protobuf:"bytes,3,opt,name=mountPath"`
}

type Hpa struct {
	// Enable is used to define whether HPA are enabled or not
	// Defaults to false.
	// +optional
	Enable bool `json:"enable,omitempty" protobuf:"varint,2,opt,name=enable"`
	// minReplicas is the lower limit for the number of replicas to which the autoscaler can scale down.
	// It defaults to 1 pod.
	// +optional
	// +kubebuilder:validation:Minimum=1
	MinReplicas *int32 `json:"minReplicas,omitempty" protobuf:"varint,2,opt,name=minReplicas"`
	// maxReplicas is the upper limit for the number of replicas to which the autoscaler can scale up.
	// It cannot be less that minReplicas.
	// +optional
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=100
	MaxReplicas *int32 `json:"maxReplicas,omitempty" protobuf:"varint,3,opt,name=maxReplicas"`
	// metrics contains the specifications for which to use to calculate the
	// desired replica count (the maximum replica count across all metrics will
	// be used).  The desired replica count is calculated multiplying the
	// ratio between the target value and the current value by the current
	// number of pods.  Ergo, metrics used must decrease as the pod count is
	// increased, and vice-versa.  See the individual metric source types for
	// more information about how each type of metric must respond.
	// +optional
	// +kubebuilder:validation:MinItems=1
	Metrics []autoscaling.MetricSpec `json:"metrics,omitempty" protobuf:"bytes,4,rep,name=metrics"`
}
