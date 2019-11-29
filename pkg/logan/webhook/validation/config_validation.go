package validation

import (
	"context"
	"github.com/logancloud/logan-app-operator/pkg/logan"
	"github.com/logancloud/logan-app-operator/pkg/logan/config"
	admssionv1beta1 "k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/runtime/inject"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"strings"
)

// ConfigValidator is the struct for webhook config validation
type ConfigValidator struct {
	client            client.Client
	decoder           *admission.Decoder
	OperatorNamespace string
}

var _ admission.Handler = &ConfigValidator{}

// Handle will handle webhook config validation
func (vHandler *ConfigValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	if !vHandler.targetConfig(req) {
		return admission.ValidationResponse(true, "")
	}

	msg, valid, err := vHandler.Validate(req)
	if err != nil {
		logger.Error(err, msg)
		return admission.Errored(http.StatusBadRequest, err)
	}

	if !valid {
		return admission.ValidationResponse(false, msg)
	}

	return admission.ValidationResponse(true, "")
}

var _ inject.Client = &ConfigValidator{}

// InjectClient will inject client into ConfigValidator
func (vHandler *ConfigValidator) InjectClient(c client.Client) error {
	vHandler.client = c
	return nil
}

var _ admission.DecoderInjector = &ConfigValidator{}

// InjectDecoder will inject decoder into ConfigValidator
func (vHandler *ConfigValidator) InjectDecoder(d *admission.Decoder) error {
	vHandler.decoder = d
	return nil
}

// Validate will do the validating for request configmap.
// Returns
//   msg: Error message
//   valid: true if valid, otherwise false
//   error: decoding error, otherwise nil
func (vHandler *ConfigValidator) Validate(req admission.Request) (string, bool, error) {
	operation := req.AdmissionRequest.Operation
	if operation == admssionv1beta1.Delete {
		return "can not delete operator's configmap", false, nil
	}

	configmap, err := vHandler.decodeConfigmap(req, vHandler.decoder)
	if err != nil {
		return "Decoding request error", false, err
	}

	if configmap == nil {
		logger.Info("Can not recognize the configmap",
			"kind", req.AdmissionRequest.Kind.Kind,
			"name", req.AdmissionRequest.Name,
			"namespace", req.AdmissionRequest.Namespace)
		return "Can not decoding configmap", false, nil
	}

	text, ok := configmap.Data[logan.ConfigFilename]
	if !ok {
		logger.Info("Can not find config.yaml in the configmap",
			"kind", req.AdmissionRequest.Kind.Kind,
			"name", req.AdmissionRequest.Name,
			"namespace", req.AdmissionRequest.Namespace)
		return "Can not find config.yaml in the configmap", false, nil
	}

	if strings.TrimSpace(text) == "" {
		logger.Info("config.yaml in the configmap can not blank",
			"kind", req.AdmissionRequest.Kind.Kind,
			"name", req.AdmissionRequest.Name,
			"namespace", req.AdmissionRequest.Namespace)
		return "config.yaml in the configmap can not blank", false, nil
	}

	err = config.NewConfigFromString(text)
	if err != nil {
		return "Decoding config.yaml error", false, err
	}

	return "", true, nil
}

func (vHandler *ConfigValidator) targetConfig(req admission.Request) bool {
	if req.AdmissionRequest.Name == logan.OperConfigmap &&
		req.AdmissionRequest.Namespace == vHandler.OperatorNamespace {
		return true
	}
	return false
}

func (vHandler *ConfigValidator) decodeConfigmap(req admission.Request, decoder *admission.Decoder) (*corev1.ConfigMap, error) {
	configmap := &corev1.ConfigMap{}
	err := decoder.Decode(req, configmap)
	if err != nil {
		return nil, err
	}
	return configmap, nil
}
