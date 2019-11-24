package framework

import (
	"context"
	"github.com/logancloud/logan-app-operator/pkg/logan/util/keys"
	"github.com/onsi/gomega"
	scheduling "k8s.io/api/scheduling/v1beta1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"log"
)

// SamplePriority will return specific Priority object according to boot key
func SamplePriority(bootKey types.NamespacedName) *scheduling.PriorityClass {
	return SamplePriorityWithName(bootKey, bootKey.Name)
}

// SamplePriorityWithName will return specific Priority object according to boot key and name
func SamplePriorityWithName(bootKey types.NamespacedName, name string) *scheduling.PriorityClass {
	return &scheduling.PriorityClass{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
			Annotations: map[string]string{
				keys.BootPriorityAnnotaionKeyPrefix + bootKey.Namespace: "true",
			},
		},
		Value: 10000,
	}
}

// CreatePriority will create specific Priority object
func CreatePriority(obj runtime.Object) {
	err := framework.Mgr.GetClient().Create(context.TODO(), obj)
	if apierrors.IsInvalid(err) {
		log.Printf("failed to create object, got an invalid object error: %s", err.Error())
		return
	}
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
	WaitDefaultUpdate()
}

// DeletePriority will delete specific Priority object
func DeletePriority(obj runtime.Object) {
	err := framework.Mgr.GetClient().Delete(context.TODO(), obj)
	if err != nil {
		log.Printf("failed to Delete Priority, got an invalid object error: %s", err.Error())
		return
	}
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
}
