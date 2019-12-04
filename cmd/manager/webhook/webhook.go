package webhook

import (
	"github.com/go-logr/logr"
	bootmutation "github.com/logancloud/logan-app-operator/pkg/logan/webhook/mutation"
	bootvalidation "github.com/logancloud/logan-app-operator/pkg/logan/webhook/validation"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

const (
	defaultPort        = 8443
	bootMutatorPath    = "/boot-mutator"
	bootValidatorPath  = "/boot-validator"
	bootConfigmapsPath = "/boot-configmaps"
)

// RegisterWebhook will register webhook for mutation and validation
func RegisterWebhook(mgr manager.Manager, log logr.Logger, operatorNs string) {
	hookServer := mgr.GetWebhookServer()
	hookServer.Port = defaultPort
	hookServer.Register(bootMutatorPath, &webhook.Admission{
		Handler: &bootmutation.BootMutator{
			Schema:   mgr.GetScheme(),
			Recorder: mgr.GetEventRecorderFor("logan-webhook-mutation"),
		},
	})

	hookServer.Register(bootValidatorPath, &webhook.Admission{
		Handler: &bootvalidation.BootValidator{
			Schema:   mgr.GetScheme(),
			Recorder: mgr.GetEventRecorderFor("logan-webhook-validation"),
		},
	})

	hookServer.Register(bootConfigmapsPath, &webhook.Admission{
		Handler: &bootvalidation.ConfigValidator{
			OperatorNamespace: operatorNs,
		},
	})
}
