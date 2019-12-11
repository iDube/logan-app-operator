package e2e

import (
	bootv1 "github.com/logancloud/logan-app-operator/pkg/apis/app/v1"
	operatorFramework "github.com/logancloud/logan-app-operator/test/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	autoscaling "k8s.io/api/autoscaling/v2beta1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("Testing Hpa ", func() {
	var bootKey types.NamespacedName
	var javaBoot *bootv1.JavaBoot

	BeforeEach(func() {
		// Gen new boot
		bootKey = operatorFramework.GenResource()
		bootKey.Namespace = namespace
		javaBoot = operatorFramework.SampleBoot(bootKey)
	})

	AfterEach(func() {
		// Clean boot
		operatorFramework.DeleteBootIgnoreError(javaBoot)
	})

	Describe("testing Hpa", func() {
		minReplicas := int32(1)
		maxReplicas := int32(2)
		targetAverageUtilization := int32(70)
		Context("test create boot default use hpa", func() {
			It("testing create boot default use hpa", func() {
				javaBoot.Spec.Hpa = &bootv1.Hpa{
					Enable:      true,
					MinReplicas: &minReplicas,
					MaxReplicas: &maxReplicas,
					Metrics: []autoscaling.MetricSpec{
						{
							Type: autoscaling.ResourceMetricSourceType,
							Resource: &autoscaling.ResourceMetricSource{
								Name:                     corev1.ResourceCPU,
								TargetAverageUtilization: &targetAverageUtilization,
							},
						}},
				}
				operatorFramework.CreateBoot(javaBoot)

				boot := operatorFramework.GetBoot(bootKey)
				Expect(boot.Name).Should(Equal(bootKey.Name))

				hpa := operatorFramework.GetHorizontalPodAutoscaler(bootKey)
				Expect(hpa.Name).Should(Equal(bootKey.Name))
				Expect(hpa.Spec.ScaleTargetRef.Name).Should(Equal(bootKey.Name))
				Expect(hpa.Spec.MaxReplicas).Should(Equal(maxReplicas))
				Expect(*hpa.Spec.MinReplicas).Should(Equal(minReplicas))
				Expect(*hpa.Spec.Metrics[0].Resource.TargetAverageUtilization).Should(Equal(targetAverageUtilization))
				Expect(hpa.Spec.Metrics[0].Resource.Name).Should(Equal(corev1.ResourceCPU))
				operatorFramework.DeleteHorizontalPodAutoscaler(hpa)
				operatorFramework.WaitUpdate(2)
				hpa = operatorFramework.GetHorizontalPodAutoscaler(bootKey)
				Expect(hpa.Name).Should(Equal(bootKey.Name))
			})
		})

		Context("test update boot's hpa", func() {
			It("testing update boot's hpa", func() {
				newMinReplicas := int32(2)
				newMaxReplicas := int32(3)
				newTargetAverageUtilization := int32(71)
				e2e := &operatorFramework.E2E{
					Build: func() {
						javaBoot.Spec.Hpa = &bootv1.Hpa{
							Enable:      true,
							MinReplicas: &minReplicas,
							MaxReplicas: &maxReplicas,
							Metrics: []autoscaling.MetricSpec{
								{
									Type: autoscaling.ResourceMetricSourceType,
									Resource: &autoscaling.ResourceMetricSource{
										Name:                     corev1.ResourceCPU,
										TargetAverageUtilization: &targetAverageUtilization,
									},
								}},
						}
						operatorFramework.CreateBoot(javaBoot)
					},
					Check: func() {
						boot := operatorFramework.GetBoot(bootKey)
						Expect(boot.Name).Should(Equal(bootKey.Name))

						hpa := operatorFramework.GetHorizontalPodAutoscaler(bootKey)
						Expect(hpa.Name).Should(Equal(bootKey.Name))
						Expect(hpa.Spec.ScaleTargetRef.Name).Should(Equal(bootKey.Name))
						Expect(hpa.Spec.MaxReplicas).Should(Equal(maxReplicas))
						Expect(*hpa.Spec.MinReplicas).Should(Equal(minReplicas))
						Expect(*hpa.Spec.Metrics[0].Resource.TargetAverageUtilization).Should(Equal(targetAverageUtilization))
						Expect(hpa.Spec.Metrics[0].Resource.Name).Should(Equal(corev1.ResourceCPU))
					},
					Update: func() {
						boot := operatorFramework.GetBoot(bootKey)
						boot.Spec.Hpa = &bootv1.Hpa{
							Enable:      true,
							MinReplicas: &newMinReplicas,
							MaxReplicas: &newMaxReplicas,
							Metrics: []autoscaling.MetricSpec{
								{
									Type: autoscaling.ResourceMetricSourceType,
									Resource: &autoscaling.ResourceMetricSource{
										Name:                     corev1.ResourceMemory,
										TargetAverageUtilization: &newTargetAverageUtilization,
									},
								}},
						}
						operatorFramework.UpdateBoot(boot)
					},
					Recheck: func() {
						boot := operatorFramework.GetBoot(bootKey)
						Expect(boot.Name).Should(Equal(bootKey.Name))

						hpa := operatorFramework.GetHorizontalPodAutoscaler(bootKey)
						Expect(hpa.Name).Should(Equal(bootKey.Name))
						Expect(hpa.Spec.ScaleTargetRef.Name).Should(Equal(bootKey.Name))
						Expect(hpa.Spec.MaxReplicas).Should(Equal(newMaxReplicas))
						Expect(*hpa.Spec.MinReplicas).Should(Equal(newMinReplicas))
						Expect(*hpa.Spec.Metrics[0].Resource.TargetAverageUtilization).Should(Equal(newTargetAverageUtilization))
						Expect(hpa.Spec.Metrics[0].Resource.Name).Should(Equal(corev1.ResourceMemory))
					},
				}

				e2e.Run()
			})
		})
	})
})
