package operator

import (
	"context"
	"fmt"
	v1 "github.com/logancloud/logan-app-operator/pkg/apis/app/v1"
	loganMetrics "github.com/logancloud/logan-app-operator/pkg/logan/metrics"
	"github.com/logancloud/logan-app-operator/pkg/logan/util"
	"github.com/logancloud/logan-app-operator/pkg/logan/util/keys"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// reconcileWorkloadCreate handle create logic for workload
func (handler *BootHandler) reconcileWorkloadCreate() (*corev1.PodTemplateSpec, reconcile.Result, bool, error) {
	workload := handler.Boot.Spec.Workload
	if workload == v1.Deployment || string(workload) == "" {
		return handler.reconcileDeploymentCreate()
	} else if workload == v1.StatefulSet {
		return handler.reconcileStatefulSetCreate()
	}

	//should not execute this
	return nil, reconcile.Result{Requeue: true}, true, nil
}

// reconcileDeploymentCreate handle create logic for Deployment
func (handler *BootHandler) reconcileDeploymentCreate() (*corev1.PodTemplateSpec, reconcile.Result, bool, error) {
	boot := handler.Boot
	logger := handler.Logger
	c := handler.Client
	requeue := false

	depFound := &appsv1.Deployment{}
	workloadName := WorkloadName(boot)
	err := c.Get(context.TODO(), types.NamespacedName{Name: workloadName, Namespace: boot.Namespace}, depFound)
	if err != nil {
		if errors.IsNotFound(err) {
			dep := handler.NewDeployment()
			logger.Info("Creating Deployment", "deploy containers", dep.Spec.Template.Spec.Containers)
			err = c.Create(context.TODO(), dep)
			if err != nil {
				msg := fmt.Sprintf("Failed to create Deployment: %s", workloadName)
				logger.Error(err, msg)
				loganMetrics.UpdateReconcileErrors(boot.Kind,
					loganMetrics.RECONCILE_CREATE_STAGE,
					loganMetrics.RECONCILE_CREATE_DEPLOYMENT_SUBSTAGE,
					boot.Name)
				handler.RecordEvent(keys.FailedCreateDeployment, msg, err)
				return nil, reconcile.Result{}, true, err
			}

			handler.RecordEvent(keys.CreatedDeployment, fmt.Sprintf("Created Deployment: %s", workloadName), nil)
			depFound = dep
			requeue = true
		} else {
			msg := fmt.Sprintf("Failed to get Deployment: %s", workloadName)
			logger.Error(err, msg)
			loganMetrics.UpdateReconcileErrors(boot.Kind,
				loganMetrics.RECONCILE_CREATE_STAGE,
				loganMetrics.RECONCILE_GET_DEPLOYMENT_SUBSTAGE,
				boot.Name)
			handler.RecordEvent(keys.FailedGetDeployment, msg, err)
			return nil, reconcile.Result{}, true, err
		}
	}

	return &depFound.Spec.Template, reconcile.Result{Requeue: requeue}, requeue, err
}

// reconcileStatefulSetCreate handle create logic for StatefulSet
func (handler *BootHandler) reconcileStatefulSetCreate() (*corev1.PodTemplateSpec, reconcile.Result, bool, error) {
	boot := handler.Boot
	logger := handler.Logger
	c := handler.Client
	requeue := false

	stsFound := &appsv1.StatefulSet{}
	workloadName := WorkloadName(boot)
	err := c.Get(context.TODO(), types.NamespacedName{Name: workloadName, Namespace: boot.Namespace}, stsFound)
	if err != nil {
		if errors.IsNotFound(err) {
			sts := handler.NewStatefulSet()
			logger.Info("Creating StatefulSet", "statefulSet containers", sts.Spec.Template.Spec.Containers)
			err = c.Create(context.TODO(), sts)
			if err != nil {
				msg := fmt.Sprintf("Failed to create statefulSet: %s", workloadName)
				logger.Error(err, msg)
				loganMetrics.UpdateReconcileErrors(boot.Kind,
					loganMetrics.RECONCILE_CREATE_STAGE,
					loganMetrics.RECONCILE_CREATE_STATEFULSET_SUBSTAGE,
					boot.Name)
				handler.RecordEvent(keys.FailedCreateStatefulSet, msg, err)
				return nil, reconcile.Result{}, true, err
			}

			handler.RecordEvent(keys.CreatedStatefulSet, fmt.Sprintf("Created StatefulSet: %s", workloadName), nil)
			stsFound = sts
			requeue = true
		} else {
			msg := fmt.Sprintf("Failed to get StatefulSet: %s", workloadName)
			logger.Error(err, msg)
			loganMetrics.UpdateReconcileErrors(boot.Kind,
				loganMetrics.RECONCILE_CREATE_STAGE,
				loganMetrics.RECONCILE_GET_STATEFULSET_SUBSTAGE,
				boot.Name)
			handler.RecordEvent(keys.FailedGetStatefulSet, msg, err)
			return nil, reconcile.Result{}, true, err
		}
	}

	return &stsFound.Spec.Template, reconcile.Result{Requeue: requeue}, requeue, err
}

// reconcileWorkloadUpdate handle update logic for workload
func (handler *BootHandler) reconcileWorkloadUpdate() (*corev1.PodTemplateSpec, reconcile.Result, bool, error) {
	workload := handler.Boot.Spec.Workload
	if workload == v1.Deployment || string(workload) == "" {
		return handler.reconcileDeploymentUpdate()
	} else if workload == v1.StatefulSet {
		return handler.reconcileStatefulSetUpdate()
	}

	//should not execute this
	return nil, reconcile.Result{Requeue: true}, true, nil
}

// reconcileDeploymentUpdate handle update logic for deployment
func (handler *BootHandler) reconcileDeploymentUpdate() (*corev1.PodTemplateSpec, reconcile.Result, bool, error) {
	boot := handler.Boot
	logger := handler.Logger
	c := handler.Client

	depFound := &appsv1.Deployment{}
	workloadName := WorkloadName(boot)
	err := c.Get(context.TODO(), types.NamespacedName{Name: workloadName, Namespace: boot.Namespace}, depFound)
	if err != nil {
		logger.Error(err, "Failed to get Deployment")
		loganMetrics.UpdateReconcileErrors(boot.Kind,
			loganMetrics.RECONCILE_UPDATE_STAGE,
			loganMetrics.RECONCILE_GET_DEPLOYMENT_SUBSTAGE,
			boot.Name)
		return nil, reconcile.Result{Requeue: true}, true, err
	}
	result, requeue, err := handler.innerReconcileUpdateDeploy(depFound)
	return &depFound.Spec.Template, result, requeue, err
}

// innerReconcileUpdateDeploy handle update logic of Deployment
func (handler *BootHandler) innerReconcileUpdateDeploy(deploy *appsv1.Deployment) (reconcile.Result, bool, error) {
	logger := handler.Logger
	boot := handler.Boot
	c := handler.Client

	updated := false
	rebootUpdated := false
	restartUpdated := false

	reason := "Updating Deployment"
	// 1. Check ownerReferences
	ownerReferences := deploy.OwnerReferences
	if ownerReferences == nil || len(ownerReferences) == 0 {
		logger.Info(reason, "type", "ownerReferences", "deploy", deploy.Name)

		_ = controllerutil.SetControllerReference(handler.OperatorBoot, deploy, handler.Scheme)

		updated = true
	}

	// 2. Check size
	size := boot.Spec.Replicas
	if *deploy.Spec.Replicas != *size {
		logger.Info(reason, "type", "replicas", "deploy", deploy.Name,
			"old", deploy.Spec.Replicas, "new", size)
		*deploy.Spec.Replicas = *size

		updated = true
	}

	// 3. Check pod spec
	restartUpdated, rebootUpdated, err := handler.reconcilePodTemplateSpecUpdate(&deploy.Spec.Template)
	if err != nil {
		return reconcile.Result{Requeue: true}, true, err
	}

	if rebootUpdated {
		updateDeploy := handler.NewDeployment()
		deploy.Spec = updateDeploy.Spec
		logger.Info("this update will cause rolling update", "Deploy", deploy.Name)
	}

	if updated || rebootUpdated || restartUpdated {
		err := c.Update(context.TODO(), deploy)
		if err != nil {
			msg := fmt.Sprintf("Failed to update Deployment: %s", deploy.GetName())
			logger.Info(msg, "err", err.Error())
			loganMetrics.UpdateReconcileErrors(boot.Kind,
				loganMetrics.RECONCILE_UPDATE_STAGE,
				loganMetrics.RECONCILE_UPDATE_DEPLOYMENT_SUBSTAGE,
				boot.Name)
			handler.RecordEvent(keys.FailedUpdateDeployment, msg, err)

			return reconcile.Result{Requeue: true}, true, err
		}

		handler.RecordEvent(keys.UpdatedDeployment, fmt.Sprintf("Updated Deployment: %s", deploy.GetName()), nil)
		return reconcile.Result{Requeue: true}, true, nil
	}

	return reconcile.Result{}, false, nil
}

// reconcileStatefulSetUpdate handle update logic for StatefulSet
func (handler *BootHandler) reconcileStatefulSetUpdate() (*corev1.PodTemplateSpec, reconcile.Result, bool, error) {
	boot := handler.Boot
	logger := handler.Logger
	c := handler.Client

	stsFound := &appsv1.StatefulSet{}
	workloadName := WorkloadName(boot)

	err := c.Get(context.TODO(), types.NamespacedName{Name: workloadName, Namespace: boot.Namespace}, stsFound)
	if err != nil {
		logger.Error(err, "Failed to get StatefulSet")
		loganMetrics.UpdateReconcileErrors(boot.Kind,
			loganMetrics.RECONCILE_UPDATE_STAGE,
			loganMetrics.RECONCILE_GET_STATEFULSET_SUBSTAGE,
			boot.Name)
		return nil, reconcile.Result{Requeue: true}, true, err
	}
	result, requeue, err := handler.innerReconcileUpdateStatefulSet(stsFound)
	return &stsFound.Spec.Template, result, requeue, err
}

// innerReconcileUpdateDeploy handle update logic of StatefulSet
func (handler *BootHandler) innerReconcileUpdateStatefulSet(sts *appsv1.StatefulSet) (reconcile.Result, bool, error) {
	logger := handler.Logger
	boot := handler.Boot
	c := handler.Client

	updated := false
	rebootUpdated := false
	restartUpdated := false

	reason := "Updating StatefulSet"
	// 1. Check ownerReferences
	ownerReferences := sts.OwnerReferences
	if ownerReferences == nil || len(ownerReferences) == 0 {
		logger.Info(reason, "type", "ownerReferences", "statefulset", sts.Name)

		_ = controllerutil.SetControllerReference(handler.OperatorBoot, sts, handler.Scheme)

		updated = true
	}

	// 2. Check size
	size := boot.Spec.Replicas
	if *sts.Spec.Replicas != *size {
		logger.Info(reason, "type", "replicas", "statefulset", sts.Name,
			"old", sts.Spec.Replicas, "new", size)
		*sts.Spec.Replicas = *size

		updated = true
	}

	// 3. Check pod spec
	restartUpdated, rebootUpdated, err := handler.reconcilePodTemplateSpecUpdate(&sts.Spec.Template)
	if err != nil {
		return reconcile.Result{Requeue: true}, true, err
	}

	if rebootUpdated {
		updateSts := handler.NewStatefulSet()
		sts.Spec = updateSts.Spec
		logger.Info("this update will cause rolling update", "statefulset", sts.Name)
	}

	if updated || rebootUpdated || restartUpdated {
		err := c.Update(context.TODO(), sts)
		if err != nil {
			msg := fmt.Sprintf("Failed to update StatefulSet: %s", sts.GetName())
			logger.Info(msg, "err", err.Error())
			loganMetrics.UpdateReconcileErrors(boot.Kind,
				loganMetrics.RECONCILE_UPDATE_STAGE,
				loganMetrics.RECONCILE_UPDATE_STATEFULSET_SUBSTAGE,
				boot.Name)
			handler.RecordEvent(keys.FailedUpdateStatefulSet, msg, err)

			return reconcile.Result{Requeue: true}, true, err
		}

		handler.RecordEvent(keys.UpdatedStatefulSet, fmt.Sprintf("Updated StatefulSet: %s", sts.GetName()), nil)
		return reconcile.Result{Requeue: true}, true, nil
	}

	return reconcile.Result{}, false, nil
}

// reconcilePodTemplateSpecUpdate handle update logic of PodTemplateSpec
func (handler *BootHandler) reconcilePodTemplateSpecUpdate(podSpec *corev1.PodTemplateSpec) (bool, bool, error) {
	boot := handler.Boot
	workloadName := WorkloadName(boot)
	logger := handler.Logger.WithValues("Workload", workloadName)
	rebootUpdated := false
	restartUpdated := false
	reason := "Updating Workload PodTemplateSpec"
	// "spec.template.spec.containers" is a required value, no need to verify.
	// 1. Check image and version:
	workloadImg := podSpec.Spec.Containers[0].Image
	bootImg := AppContainerImageName(handler.Boot, handler.Config.AppSpec)
	if bootImg != workloadImg {
		logger.Info(reason, "type", "image",
			"old", workloadImg, "new", bootImg)

		rebootUpdated = true
	}

	// 2. Check env: check fist container(boot container)
	workloadEnv := podSpec.Spec.Containers[0].Env
	bootEnv := boot.Spec.Env
	if !reflect.DeepEqual(workloadEnv, bootEnv) {
		logger.Info(reason, "type", "env",
			"old", workloadEnv, "new", bootEnv)

		rebootUpdated = true
	}

	// 3. Check port: check fist container(boot container)
	workloadPorts := podSpec.Spec.Containers[0].Ports
	bootPorts := []corev1.ContainerPort{{
		Name:          HttpPortName,
		ContainerPort: boot.Spec.Port,
		Protocol:      corev1.ProtocolTCP}}
	if !reflect.DeepEqual(workloadPorts, bootPorts) {
		logger.Info(reason, "type", "port",
			"old", workloadPorts, "new", bootPorts)

		rebootUpdated = true
	}

	// 4. Check resources: check fist container(boot container)
	workloadResources := podSpec.Spec.Containers[0].Resources
	bootResources := boot.Spec.Resources
	if !reflect.DeepEqual(workloadResources, bootResources) {
		logger.Info(reason, "type", "resources",
			"old", workloadResources, "new", bootResources)

		rebootUpdated = true
	}

	// 5. Check liveness and readiness : check fist container(boot container)
	// 5.1 Check liveness
	livenessProbe := podSpec.Spec.Containers[0].LivenessProbe
	bootHealth := *boot.Spec.Health
	if bootHealth == "" {
		if livenessProbe != nil {
			// Remove the 2 existing probes.
			workloadHealth := livenessProbe.HTTPGet.Path
			logger.Info(reason, "type", "health",
				"old", workloadHealth, "new", "")

			rebootUpdated = true
		}
	} else {
		if livenessProbe == nil {
			// 1. If probe is nil, add Liveness and Readiness
			logger.Info(reason, "type", "health",
				"old", "empty", "new", bootHealth)

			rebootUpdated = true
		} else {
			workloadHealth := livenessProbe.HTTPGet.Path
			// 2. If livenessProbe is not nil, we only need to update the health path
			if workloadHealth != bootHealth {
				logger.Info(reason, "type", "health",
					"old", workloadHealth, "new", bootHealth)

				rebootUpdated = true
			}
		}
	}

	// 5.2 Check readiness
	readinessProbe := podSpec.Spec.Containers[0].ReadinessProbe
	readinessPath := *boot.Spec.Health

	// if boot.Spec.Health is empty, ignore the Readiness
	if boot.Spec.Readiness != nil && *boot.Spec.Readiness != "" &&
		boot.Spec.Health != nil && *boot.Spec.Health != "" {
		readinessPath = *boot.Spec.Readiness
	}

	if readinessPath == "" {
		if readinessProbe != nil {
			// Remove the 2 existing probes.
			workloadHealth := readinessProbe.HTTPGet.Path
			logger.Info(reason, "type", "readiness",
				"old", workloadHealth, "new", "")

			rebootUpdated = true
		}
	} else {
		if readinessProbe == nil {
			// 1. If probe is nil, add  Readiness
			logger.Info(reason, "type", "readiness",
				"old", "empty", "new", readinessPath)

			rebootUpdated = true
		} else {
			workloadHealth := readinessProbe.HTTPGet.Path
			// 2. If readinessProbe is not nil, we only need to update the readiness path
			if workloadHealth != readinessPath {
				logger.Info(reason, "type", "readiness",
					"old", workloadHealth, "new", readinessPath)

				rebootUpdated = true
			}
		}
	}

	// 6. Check nodeSelector: map[string]string
	workloadNodeSelector := podSpec.Spec.NodeSelector
	bootNodeSelector := boot.Spec.NodeSelector
	if !reflect.DeepEqual(workloadNodeSelector, bootNodeSelector) {
		logger.Info(reason, "type", "nodeSelector",
			"old", workloadNodeSelector, "new", bootNodeSelector)

		rebootUpdated = true
	}

	// 7. Check command
	workloadCommand := podSpec.Spec.Containers[0].Command
	bootCommand := boot.Spec.Command
	if !reflect.DeepEqual(workloadCommand, bootCommand) {
		logger.Info(reason, "type", "command",
			"old", workloadCommand, "new", bootCommand)

		rebootUpdated = true
	}

	// 8. Check vol
	workloadVols := podSpec.Spec.Containers[0].VolumeMounts
	bootVolStr, ok := boot.Annotations[keys.BootDeployPvcsAnnotationKey]
	if ok && bootVolStr != "" {
		bootVols, err := DecodeVolumeMountVars(bootVolStr)
		if err != nil {
			logger.Error(err, "can not decode VolumeMount",
				keys.BootDeployPvcsAnnotationKey, bootVolStr)
			return true, true, err
		}

		if !VolumeMountVarsEq(workloadVols, bootVols) {
			deleted, added, modified := util.DifferenceVol(workloadVols, bootVols)
			logger.Info("Boot VolumeMounts change.",
				"deleted", deleted, "added", added, "modified", modified)

			volUpdated, err := handler.checkVolumeMountUpdate(deleted, added, modified)
			if err != nil {
				logger.Error(err, "Fail to reconcile VolumeMounts",
					"deleted", deleted, "added", added, "modified", modified)
				return true, true, err
			}

			if volUpdated {
				logger.Info(reason, "type", "VolumeMounts",
					"old", workloadVols, "new", bootVols, keys.BootDeployPvcsAnnotationKey, bootVolStr)
				rebootUpdated = true
			}
		}
	} else if workloadVols != nil {
		// if we have this, is very surprised
		logger.Info(reason, "type", "VolumeMounts",
			"old", workloadVols, "new", nil)
		// rebootUpdated = true
	}

	// 9. Check RestartedAt
	if bootRestartedAt, ok := boot.Annotations[keys.BootRestartedAtAnnotationKey]; ok {
		if podSpec.Annotations == nil {
			podSpec.Annotations = make(map[string]string)
		}

		workloadRestartedAt, ok := podSpec.Annotations[keys.BootRestartedAtAnnotationKey]
		if !ok || bootRestartedAt != workloadRestartedAt {
			podSpec.Annotations[keys.BootRestartedAtAnnotationKey] = bootRestartedAt
			restartUpdated = true
		}
	}

	// 10. Check Priority
	if podSpec.Spec.PriorityClassName != boot.Spec.Priority {
		logger.Info(reason, "type", "Priority",
			"old", podSpec.Spec.PriorityClassName, "new", boot.Spec.Priority)
		rebootUpdated = true
	}

	return restartUpdated, rebootUpdated, nil
}

// getWorkloadStatus will return Replicas and CurrentReplicas
func (handler *BootHandler) getWorkloadStatus() (int32, int32, error) {
	logger := handler.Logger
	boot := handler.Boot
	c := handler.Client

	workloadName := WorkloadName(boot)
	if boot.Spec.Workload == v1.Deployment {
		dep := &appsv1.Deployment{}
		err := c.Get(context.TODO(), types.NamespacedName{Name: workloadName, Namespace: boot.Namespace}, dep)
		if err != nil {
			logger.Error(err, "Failed to get Deployment")
			loganMetrics.UpdateReconcileErrors(boot.Kind,
				loganMetrics.RECONCILE_UPDATE_BOOT_META_STAGE,
				loganMetrics.RECONCILE_GET_DEPLOYMENT_SUBSTAGE,
				boot.Name)
			return 0, 0, err
		}
		return dep.Status.Replicas, dep.Status.AvailableReplicas, nil
	} else if boot.Spec.Workload == v1.StatefulSet {
		sts := &appsv1.StatefulSet{}
		err := c.Get(context.TODO(), types.NamespacedName{Name: workloadName, Namespace: boot.Namespace}, sts)
		if err != nil {
			logger.Error(err, "Failed to get StatefulSet")
			loganMetrics.UpdateReconcileErrors(boot.Kind,
				loganMetrics.RECONCILE_UPDATE_BOOT_META_STAGE,
				loganMetrics.RECONCILE_GET_STATEFULSET_SUBSTAGE,
				boot.Name)
			return 0, 0, err
		}
		return sts.Status.Replicas, sts.Status.CurrentReplicas, nil
	}

	//should not execute this
	return 0, 0, nil
}
