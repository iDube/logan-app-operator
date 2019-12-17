package framework

import (
	"github.com/onsi/gomega"
	autoscaling "k8s.io/api/autoscaling/v2beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

// GetHorizontalPodAutoscaler will return specific HorizontalPodAutoscaler from kubernetes by NamespacedName
func GetHorizontalPodAutoscaler(nn types.NamespacedName) *autoscaling.HorizontalPodAutoscaler {
	hpa := &autoscaling.HorizontalPodAutoscaler{}
	var err error
	gomega.Eventually(func() error {
		hpa, err = framework.KubeClient.AutoscalingV2beta1().HorizontalPodAutoscalers(nn.Namespace).Get(nn.Name, metav1.GetOptions{})
		return err
	}, defaultTimeout).
		Should(gomega.Succeed())
	return hpa
}

// DeleteHorizontalPodAutoscaler will delete specific HorizontalPodAutoscaler from kubernetes
func DeleteHorizontalPodAutoscaler(hpa *autoscaling.HorizontalPodAutoscaler) {
	gomega.Eventually(func() error {
		return framework.KubeClient.AutoscalingV2beta1().HorizontalPodAutoscalers(hpa.Namespace).Delete(hpa.Name, &metav1.DeleteOptions{})
	}, defaultTimeout).
		Should(gomega.Succeed())
}
