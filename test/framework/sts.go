package framework

import (
	"github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// GetStatefulSet will return specific StatefulSet from kubernetes by NamespacedName
func GetStatefulSet(nn types.NamespacedName) *appsv1.StatefulSet {
	sts := &appsv1.StatefulSet{}
	var err error
	gomega.Eventually(func() error {
		sts, err = framework.KubeClient.AppsV1().StatefulSets(nn.Namespace).Get(nn.Name, metav1.GetOptions{})
		return err
	}, defaultTimeout).
		Should(gomega.Succeed())
	return sts
}

// DeleteStatefulSet will delete specific StatefulSet from kubernetes
func DeleteStatefulSet(sts *appsv1.StatefulSet) {
	gomega.Eventually(func() error {
		return framework.KubeClient.AppsV1().StatefulSets(sts.Namespace).Delete(sts.Name, &metav1.DeleteOptions{})
	}, defaultTimeout).
		Should(gomega.Succeed())
}
