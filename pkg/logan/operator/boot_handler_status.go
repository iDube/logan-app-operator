package operator

import (
	appv1 "github.com/logancloud/logan-app-operator/pkg/apis/app/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"strconv"
)

// ReconcileUpdateStatus handle update logic for status
func (handler *BootHandler) ReconcileUpdateStatus() (reconcile.Result, bool, bool, error) {
	logger := handler.Logger
	bootStatus := handler.OperatorStatus
	bootSpec := handler.OperatorSpec
	boot := handler.Boot
	c := handler.Client
	changed := false

	reason := "Update status"

	//status
	// 1. Replicas
	if bootStatus.Replicas != *bootSpec.Replicas {
		bootStatus.Replicas = *bootSpec.Replicas
		logger.Info(reason, "type", "status.Replicas", "to", bootStatus.Replicas)
		changed = true
	}

	if bootStatus.HPAReplicas != *bootSpec.Replicas {
		bootStatus.HPAReplicas = *bootSpec.Replicas
		logger.Info(reason, "type", "status.HPAReplicas", "to", bootStatus.HPAReplicas)
		changed = true
	}

	readyReplicas, currentReplicas, err := handler.getWorkloadStatus()
	if err != nil {
		return reconcile.Result{Requeue: true}, true, changed, err
	}
	if bootStatus.ReadyReplicas != readyReplicas {
		bootStatus.ReadyReplicas = readyReplicas
		logger.Info(reason, "type", "status.ReadyReplicas", "to", bootStatus.ReadyReplicas)
		changed = true
	}
	if bootStatus.CurrentReplicas != currentReplicas {
		bootStatus.CurrentReplicas = currentReplicas
		logger.Info(reason, "type", "status.CurrentReplicas", "to", bootStatus.CurrentReplicas)
		changed = true
	}

	// 2. Selector
	podLabels := PodLabels(handler.Boot)
	selector := metav1.FormatLabelSelector(&metav1.LabelSelector{
		MatchLabels: podLabels,
	})

	if bootStatus.Selector == "" || bootStatus.Selector != selector {
		bootStatus.Selector = selector
		logger.Info(reason, "type", "status.Selector", "to", bootStatus.Selector)
		changed = true

	}

	// 3. Workload
	if bootSpec.Workload != "" {
		if bootStatus.Workload != bootSpec.Workload {
			bootStatus.Workload = bootSpec.Workload
			logger.Info(reason, "type", "status.Workload", "to", bootStatus.Workload)
			changed = true
		}
	}

	if bootStatus.Workload == "" {
		bootStatus.Workload = appv1.Deployment
		logger.Info(reason, "type", "status.Workload", "to", bootStatus.Workload)
		changed = true
	}

	// 4. service
	svcList, err := handler.listRuntimeService()
	if err != nil {
		return reconcile.Result{Requeue: true}, true, changed, err
	}
	svcText := TransferServiceNames(svcList.Items)
	if bootStatus.Services != svcText {
		bootStatus.Services = svcText
		logger.Info(reason, "type", "status.Services", "to", bootStatus.Services)
		changed = true
	}

	// 5. revision
	revisionLst, _ := c.ListRevision(boot.Namespace, podLabels)
	latestRevision := revisionLst.SelectLatestRevision()
	if latestRevision != nil {
		revisionId := strconv.Itoa(latestRevision.GetRevisionId())
		if bootStatus.Revision != revisionId {
			bootStatus.Revision = revisionId
			logger.Info(reason, "type", "status.Revision", "to", bootStatus.Revision)
			changed = true
		}
	}

	return reconcile.Result{Requeue: changed}, changed, changed, nil
}
