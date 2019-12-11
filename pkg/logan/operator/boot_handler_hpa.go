package operator

import (
	"context"
	"github.com/logancloud/logan-app-operator/pkg/logan/util"
	autoscaling "k8s.io/api/autoscaling/v2beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// ReconcileHpa handle update logic for HPA
func (handler *BootHandler) ReconcileHpa() (reconcile.Result, error) {
	boot := handler.Boot
	logger := handler.Logger
	c := handler.Client

	hpaFound := &autoscaling.HorizontalPodAutoscaler{}
	err := c.Get(context.TODO(), types.NamespacedName{Name: boot.Name, Namespace: boot.Namespace}, hpaFound)

	//disable hpa, should delete hpa
	if boot.Spec.Hpa == nil || boot.Spec.Hpa.Enable == false {
		if err != nil {
			if errors.IsNotFound(err) {
				// is ok
				return reconcile.Result{Requeue: false}, nil
			}
			return reconcile.Result{Requeue: true}, err
		}

		logger.Info("Exist HorizontalPodAutoscaler,need to delete.")
		err = c.Delete(context.TODO(), hpaFound)
		if err != nil {
			logger.Error(err, "Failed to delete HorizontalPodAutoscaler.")
			return reconcile.Result{Requeue: true}, err
		}
		logger.Info("Succeed to delete HorizontalPodAutoscaler.", "hpa", hpaFound.Spec)
		return reconcile.Result{Requeue: true}, nil
	} else {
		// enable hpa
		if err != nil {
			//not found,should create it
			if errors.IsNotFound(err) {
				hpa := handler.NewHpa()
				logger.Info("HorizontalPodAutoscaler not found,need to create.")
				err := c.Create(context.TODO(), hpa)
				if err != nil {
					logger.Error(err, "Failed to create HorizontalPodAutoscaler.")
					return reconcile.Result{Requeue: true}, err
				}
				logger.Info("Succeed to create HorizontalPodAutoscaler.", "hpa", hpa.Spec)
				return reconcile.Result{Requeue: true}, nil
			}
			logger.Error(err, "Failed to get HorizontalPodAutoscaler.")
			return reconcile.Result{Requeue: true}, err
		} else {
			//hpa exist, check update
			changed := false

			// 1. Check ownerReferences
			ownerReferences := hpaFound.OwnerReferences
			if ownerReferences == nil || len(ownerReferences) == 0 {
				_ = controllerutil.SetControllerReference(handler.OperatorBoot, hpaFound, handler.Scheme)
				changed = true
			}

			expectHpa := handler.NewHpa()
			if *expectHpa.Spec.MinReplicas != *hpaFound.Spec.MinReplicas {
				changed = true
			}

			if expectHpa.Spec.MaxReplicas != hpaFound.Spec.MaxReplicas {
				changed = true
			}

			deleted, added, modified := util.DifferenceMetric(expectHpa.Spec.Metrics, hpaFound.Spec.Metrics)
			if len(deleted)+len(added)+len(modified) > 0 {
				changed = true
			}

			if changed {
				logger.Info("HorizontalPodAutoscaler is too old, need to update",
					"old", hpaFound.Spec, "new", expectHpa.Spec)
				hpaFound.Spec = expectHpa.Spec
				err := c.Update(context.TODO(), hpaFound)
				if err != nil {
					logger.Error(err, "Failed to update HorizontalPodAutoscaler.")
					return reconcile.Result{Requeue: true}, err
				}
				logger.Info("Succeed to update HorizontalPodAutoscaler.")
				return reconcile.Result{Requeue: true}, nil
			}

			// is ok
			return reconcile.Result{Requeue: false}, nil
		}
	}
}

// NewHpa will create a new HPA
func (handler *BootHandler) NewHpa() *autoscaling.HorizontalPodAutoscaler {
	boot := handler.Boot

	hpa := &autoscaling.HorizontalPodAutoscaler{
		TypeMeta: metav1.TypeMeta{
			Kind:       "HorizontalPodAutoscaler",
			APIVersion: "autoscaling/v2beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      boot.Name,
			Namespace: boot.Namespace,
		},
		Spec: autoscaling.HorizontalPodAutoscalerSpec{
			MinReplicas: boot.Spec.Hpa.MinReplicas,
			MaxReplicas: *boot.Spec.Hpa.MaxReplicas,
			Metrics:     boot.Spec.Hpa.Metrics,
			ScaleTargetRef: autoscaling.CrossVersionObjectReference{
				APIVersion: boot.APIVersion,
				Kind:       boot.Kind,
				Name:       boot.Name,
			},
		},
	}

	_ = controllerutil.SetControllerReference(handler.OperatorBoot, hpa, handler.Scheme)

	return hpa
}
