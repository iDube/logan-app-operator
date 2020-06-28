package framework

import (
	"github.com/logancloud/logan-app-operator/pkg/apis"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/runtime/signals"
	"time"
)

var defaultTimeout = 1 * time.Minute

// Framework is the struct for e2e test framework
type Framework struct {
	Mgr            manager.Manager
	KubeClient     kubernetes.Interface
	HTTPClient     *http.Client
	MasterHost     string
	DefaultTimeout time.Duration
}

var (
	framework *Framework
)

// InitFramework will return framework object
func InitFramework() (*Framework, error) {
	var err error
	framework, err = New()
	return framework, err
}

// New will new the item defined in Framework structure
func New() (*Framework, error) {
	kubeconfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}

	cli, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "creating new kube-client failed")
	}

	httpc := cli.CoreV1().RESTClient().(*rest.RESTClient).Client
	if err != nil {
		return nil, errors.Wrap(err, "creating http-client failed")
	}
	mgr, err := manager.New(kubeconfig, manager.Options{
		Namespace: "",
	})

	if err != nil {
		return nil, errors.Wrap(err, "creating new manager failed")
	}

	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		return nil, errors.Wrap(err, "creating add to scheme failed")
	}

	f := &Framework{
		Mgr:            mgr,
		MasterHost:     kubeconfig.Host,
		KubeClient:     cli,
		HTTPClient:     httpc,
		DefaultTimeout: time.Minute,
	}

	// start go client list-watch
	go mgr.Start(signals.SetupSignalHandler())

	return f, nil
}
