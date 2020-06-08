package main

import (
	"context"
	"flag"
	v1 "github.com/logancloud/logan-app-operator/pkg/apis/app/v1"
	"github.com/logancloud/logan-app-operator/pkg/logan"
	"github.com/logancloud/logan-app-operator/pkg/logan/config"
	"github.com/logancloud/logan-app-operator/pkg/logan/util/keys"
	operatorFramework "github.com/logancloud/logan-app-operator/test/framework"
	"github.com/operator-framework/operator-sdk/pkg/log/zap"
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("recover")
var namespace string

func main() {
	// Add the zap logger flag set to the CLI. The flag set must
	// be added before calling pflag.Parse().
	pflag.CommandLine.AddFlagSet(zap.FlagSet())
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.StringVarP(&namespace, "namespace", "n", "", "the recover namespace")
	pflag.Parse()
	logf.SetLogger(zap.Logger())

	if namespace == "" {
		log.Info("namespace can not be empty")
		os.Exit(1)
	}
	log.Info("the env", "BIZ_ENVS", logan.BizEnvs)
	log.Info("start recover", "namespace", namespace)

	framework, err := operatorFramework.InitFramework()
	if err != nil {
		log.Error(err, "failed to setup framework: %v\n")
		os.Exit(1)
	}
	targetNamespace, err := framework.KubeClient.CoreV1().Namespaces().Get(context.TODO(), namespace, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			log.Error(err, "can not find namesapce", "namespace", namespace)
		}
		os.Exit(1)
	}

	log.Info("get target namespace.", "namespace", targetNamespace)
	revisionList := &v1.BootRevisionList{}
	err = framework.Mgr.GetClient().List(context.TODO(), revisionList,
		client.InNamespace(namespace))
	if err != nil {
		log.Error(err, "failed to get revision list")
		os.Exit(1)
	}
	log.Info("get revision list", "size", len(revisionList.Items))

	for index, revision := range revisionList.Items {
		newRevision := &v1.BootRevision{
			ObjectMeta: metav1.ObjectMeta{
				Name:        revision.Labels[keys.BootNameKey],
				Namespace:   revision.Namespace,
				Annotations: initRevisionAnnotations(&revision),
			},
			Spec:     revision.Spec,
			BootType: revision.BootType,
			AppKey:   revision.AppKey,
		}
		newHash := newRevision.BootHash()
		log.Info("process", "index", index, "revision", revision.Name, "hash", newHash)
		revision.Annotations[keys.BootRevisionHashAnnotationKey] = newHash
		framework.Mgr.GetClient().Update(context.TODO(), &revision)
	}
	log.Info("recover work done.")
}

func initRevisionAnnotations(revision *v1.BootRevision) map[string]string {
	if revision.Annotations != nil {
		if val, isok := revision.Annotations[config.BootProfileAnnotationKey]; isok {
			return map[string]string{
				config.BootProfileAnnotationKey: val,
			}
		}
	}
	return map[string]string{}
}
