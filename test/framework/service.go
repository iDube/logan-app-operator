package framework

import (
	"context"
	"github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"log"
)

// CreateService will create service in kubernetes
func CreateService(svr *corev1.Service) *corev1.Service {
	service := &corev1.Service{}
	var err error
	gomega.Eventually(func() error {
		service, err = framework.KubeClient.CoreV1().Services(svr.Namespace).
			Create(context.TODO(), svr, metav1.CreateOptions{})
		return err
	}, defaultTimeout).
		Should(gomega.Succeed())
	WaitDefaultUpdate()
	return service
}

// CreateServiceWithError will create service in kubernetes, return error if occur
func CreateServiceWithError(svr *corev1.Service) (*corev1.Service, error) {
	service := &corev1.Service{}
	var err error
	gomega.Eventually(func() error {
		service, err = framework.KubeClient.CoreV1().Services(svr.Namespace).
			Create(context.TODO(), svr, metav1.CreateOptions{})
		return err
	}, defaultTimeout).
		ShouldNot(gomega.Succeed())
	WaitDefaultUpdate()
	return service, err
}

// GetService will get service with NamespacedName from kubernetes, return service
func GetService(nn types.NamespacedName) *corev1.Service {
	service := &corev1.Service{}
	var err error
	gomega.Eventually(func() error {
		service, err = framework.KubeClient.CoreV1().Services(nn.Namespace).
			Get(context.TODO(), nn.Name, metav1.GetOptions{})
		return err
	}, defaultTimeout).
		Should(gomega.Succeed())
	WaitDefaultUpdate()
	return service
}

// GetServiceWithError will get service with NamespacedName from kubernetes, return service and error
func GetServiceWithError(nn types.NamespacedName) (*corev1.Service, error) {
	service := &corev1.Service{}
	var err error
	gomega.Eventually(func() error {
		service, err = framework.KubeClient.CoreV1().Services(nn.Namespace).
			Get(context.TODO(), nn.Name, metav1.GetOptions{})
		return err
	}, defaultTimeout).
		ShouldNot(gomega.Succeed())
	WaitDefaultUpdate()
	return service, err
}

// UpdateService will update service to kubernetes
func UpdateService(svr *corev1.Service) *corev1.Service {
	service := &corev1.Service{}
	var err error
	gomega.Eventually(func() error {
		latest := GetService(types.NamespacedName{Name: svr.Name, Namespace: svr.Namespace})
		latest.Spec = svr.Spec
		service, err = framework.KubeClient.CoreV1().Services(svr.Namespace).
			Update(context.TODO(), latest, metav1.UpdateOptions{})
		if apierrors.IsConflict(err) {
			log.Printf("failed to update object, got an Conflict error: ")
		}
		if apierrors.IsInvalid(err) {
			log.Printf("failed to update object, got an invalid object error: ")
		}
		return err
	}, defaultTimeout).
		Should(gomega.Succeed())
	WaitDefaultUpdate()
	return service
}

// DeleteService will delete service from kubernetes
func DeleteService(svr *corev1.Service) {
	gomega.Eventually(func() error {
		return framework.KubeClient.CoreV1().Services(svr.Namespace).
			Delete(context.TODO(), svr.Name, metav1.DeleteOptions{})
	}, defaultTimeout).
		Should(gomega.Succeed())
}
