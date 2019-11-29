package webhook

import (
	appv1 "github.com/logancloud/logan-app-operator/pkg/apis/app/v1"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	// ApiTypeJava is the type for JavaBoot in decoding schema
	ApiTypeJava = "JavaBoot"
	// ApiTypePhp is the type for PhpBoot in decoding schema
	ApiTypePhp = "PhpBoot"
	// ApiTypePython is the type for PythonBoot in decoding schema
	ApiTypePython = "PythonBoot"
	// ApiTypeNodeJS is the type for NodeJSBoot in decoding schema
	ApiTypeNodeJS = "NodeJSBoot"
	// ApiTypeWeb is the type for WebBoot in decoding schema
	ApiTypeWeb = "WebBoot"
)

// DecodeBoot decode the Boot object from request.
func DecodeBoot(req admission.Request, decoder *admission.Decoder) (*appv1.Boot, error) {
	bootType := req.AdmissionRequest.Kind.Kind

	var boot *appv1.Boot
	if bootType == ApiTypeJava {
		apiBoot := &appv1.JavaBoot{}
		err := decoder.Decode(req, apiBoot)
		if err != nil {
			return nil, err
		}
		boot = apiBoot.DeepCopyBoot()
	} else if bootType == ApiTypePhp {
		apiBoot := &appv1.PhpBoot{}
		err := decoder.Decode(req, apiBoot)
		if err != nil {
			return nil, err
		}
		boot = apiBoot.DeepCopyBoot()
	} else if bootType == ApiTypePython {
		apiBoot := &appv1.PythonBoot{}
		err := decoder.Decode(req, apiBoot)
		if err != nil {
			return nil, err
		}
		boot = apiBoot.DeepCopyBoot()
	} else if bootType == ApiTypeNodeJS {
		apiBoot := &appv1.NodeJSBoot{}
		err := decoder.Decode(req, apiBoot)
		if err != nil {
			return nil, err
		}
		boot = apiBoot.DeepCopyBoot()
	} else if bootType == ApiTypeWeb {
		apiBoot := &appv1.WebBoot{}
		err := decoder.Decode(req, apiBoot)
		if err != nil {
			return nil, err
		}
		boot = apiBoot.DeepCopyBoot()
	}

	return boot, nil
}

// DecodeJavaBoot decode the JavaBoot object from request.
func DecodeJavaBoot(req admission.Request, decoder *admission.Decoder) (*appv1.JavaBoot, error) {
	bootType := req.AdmissionRequest.Kind.Kind

	var boot *appv1.JavaBoot
	if bootType == ApiTypeJava {
		boot = &appv1.JavaBoot{}
		err := decoder.Decode(req, boot)
		if err != nil {
			return nil, err
		}
		return boot, nil
	}

	return boot, nil
}

// DecodePhpBoot decode the PhpBoot object from request.
func DecodePhpBoot(req admission.Request, decoder *admission.Decoder) (*appv1.PhpBoot, error) {
	bootType := req.AdmissionRequest.Kind.Kind

	var boot *appv1.PhpBoot
	if bootType == ApiTypePhp {
		boot = &appv1.PhpBoot{}
		err := decoder.Decode(req, boot)
		if err != nil {
			return nil, err
		}
		return boot, nil
	}

	return boot, nil
}

// DecodePythonBoot decode the PythonBoot object from request.
func DecodePythonBoot(req admission.Request, decoder *admission.Decoder) (*appv1.PythonBoot, error) {
	bootType := req.AdmissionRequest.Kind.Kind

	var boot *appv1.PythonBoot
	if bootType == ApiTypePython {
		boot = &appv1.PythonBoot{}
		err := decoder.Decode(req, boot)
		if err != nil {
			return nil, err
		}
		return boot, nil
	}

	return boot, nil
}

// DecodeNodeJSBoot decode the NodeJSBoot object from request.
func DecodeNodeJSBoot(req admission.Request, decoder *admission.Decoder) (*appv1.NodeJSBoot, error) {
	bootType := req.AdmissionRequest.Kind.Kind

	var boot *appv1.NodeJSBoot
	if bootType == ApiTypeNodeJS {
		boot = &appv1.NodeJSBoot{}
		err := decoder.Decode(req, boot)
		if err != nil {
			return nil, err
		}
		return boot, nil
	}

	return boot, nil
}

// DecodeWebBoot decode the WebBoot object from request.
func DecodeWebBoot(req admission.Request, decoder *admission.Decoder) (*appv1.WebBoot, error) {
	bootType := req.AdmissionRequest.Kind.Kind

	var boot *appv1.WebBoot
	if bootType == ApiTypeWeb {
		boot = &appv1.WebBoot{}
		err := decoder.Decode(req, boot)
		if err != nil {
			return nil, err
		}
		return boot, nil
	}

	return boot, nil
}
