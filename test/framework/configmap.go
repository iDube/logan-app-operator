package framework

import (
	"context"
	"github.com/onsi/gomega"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"log"
)

// GetConfigmap will return config map from kubernetes by NamespacedName
func GetConfigmap(nn types.NamespacedName) *v1.ConfigMap {
	configMap := &v1.ConfigMap{}
	var err error
	gomega.Eventually(func() error {
		configMap, err = framework.KubeClient.CoreV1().ConfigMaps(nn.Namespace).
			Get(context.TODO(), nn.Name, metav1.GetOptions{})
		return err
	}, defaultTimeout).
		Should(gomega.Succeed())
	return configMap
}

// UpdateConfigmap will update specific config map to kubernetes
func UpdateConfigmap(configMap *v1.ConfigMap) *v1.ConfigMap {
	conf := &v1.ConfigMap{}
	var err error
	gomega.Eventually(func() error {
		latest := GetConfigmap(types.NamespacedName{Name: configMap.Name, Namespace: configMap.Namespace})
		latest.Data = configMap.Data
		latest.BinaryData = configMap.BinaryData
		conf, err = framework.KubeClient.CoreV1().ConfigMaps(configMap.Namespace).
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
	return conf
}

// DeleteConfigmap will delete specific config map in kubernetes
func DeleteConfigmap(configMap *v1.ConfigMap) {
	gomega.Eventually(func() error {
		return framework.KubeClient.CoreV1().ConfigMaps(configMap.Namespace).
			Delete(context.TODO(), configMap.Name, metav1.DeleteOptions{})
	}, defaultTimeout).
		Should(gomega.Succeed())
}
