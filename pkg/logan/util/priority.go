package util

import (
	scheduling "k8s.io/api/scheduling/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// HighestUserDefinablePriority is the highest priority for user defined priority classes. Priority values larger than 1 billion are reserved for Kubernetes system use.
	// from https://github.com/kubernetes/kubernetes/blob/9016740a6ffe91bb29824f80c34087b993903bd6/pkg/apis/scheduling/types.go#L21
	HighestUserDefinablePriority = int32(1000000000)
	// SystemCriticalPriority is the beginning of the range of priority values for critical system components.
	SystemCriticalPriority = 2 * HighestUserDefinablePriority
	// SystemPriorityClassPrefix is the prefix reserved for system priority class names. Other priority
	// classes are not allowed to start with this prefix.
	SystemPriorityClassPrefix = "system-"
	// NOTE: In order to avoid conflict of names with user-defined priority classes, all the names must
	// start with SystemPriorityClassPrefix.
	SystemClusterCritical = SystemPriorityClassPrefix + "cluster-critical"
	SystemNodeCritical    = SystemPriorityClassPrefix + "node-critical"
)

// PriorityClassPermittedInNamespace returns true if we allow the given priority class name in the
// given namespace. It currently checks that system priorities are created only in the system namespace.
// from https://github.com/kubernetes/kubernetes/blob/release-1.11/plugin/pkg/admission/priority/admission.go#L144
func PriorityClassPermittedInNamespace(priorityClassName string, namespace string) bool {
	// Only allow system priorities in the system namespace. This is to prevent abuse or incorrect
	// usage of these priorities. Pods created at these priorities could preempt system critical
	// components.
	for _, spc := range SystemPriorityClasses() {
		if spc.Name == priorityClassName && namespace != metav1.NamespaceSystem {
			return false
		}
	}
	return true
}

// SystemPriorityClasses define system priority classes that are auto-created at cluster bootstrapping.
// Our API validation logic ensures that any priority class that has a system prefix or its value
// is higher than HighestUserDefinablePriority is equal to one of these SystemPriorityClasses.
// https://github.com/kubernetes/kubernetes/blob/9016740a6ffe91bb29824f80c34087b993903bd6/pkg/apis/scheduling/helpers.go#L27
var systemPriorityClasses = []*scheduling.PriorityClass{
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: SystemNodeCritical,
		},
		Value:       SystemCriticalPriority + 1000,
		Description: "Used for system critical pods that must not be moved from their current node.",
	},
	{
		ObjectMeta: metav1.ObjectMeta{
			Name: SystemClusterCritical,
		},
		Value:       SystemCriticalPriority,
		Description: "Used for system critical pods that must run in the cluster, but can be moved to another node if necessary.",
	},
}

// SystemPriorityClasses returns the list of system priority classes.
// NOTE: be careful not to modify any of elements of the returned array directly.
// from https://github.com/kubernetes/kubernetes/blob/9016740a6ffe91bb29824f80c34087b993903bd6/pkg/apis/scheduling/helpers.go#L46
func SystemPriorityClasses() []*scheduling.PriorityClass {
	return systemPriorityClasses
}
