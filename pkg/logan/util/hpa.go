package util

import (
	autoscaling "k8s.io/api/autoscaling/v2beta1"
	pathvalidation "k8s.io/apimachinery/pkg/api/validation/path"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"strings"
)

// ValidateMetrics will validate the MetricSpec
// from https://github.com/kubernetes/kubernetes/blob/9016740a6ffe91bb29824f80c34087b993903bd6/pkg/apis/autoscaling/validation/validation.go#L105
func ValidateMetrics(metrics []autoscaling.MetricSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	for i, metricSpec := range metrics {
		idxPath := fldPath.Index(i)
		if targetErrs := validateMetricSpec(metricSpec, idxPath); len(targetErrs) > 0 {
			allErrs = append(allErrs, targetErrs...)
		}
	}

	return allErrs
}

func validateCrossVersionObjectReference(ref autoscaling.CrossVersionObjectReference, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(ref.Kind) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("kind"), ""))
	} else {
		for _, msg := range pathvalidation.IsValidPathSegmentName(ref.Kind) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("kind"), ref.Kind, msg))
		}
	}

	if len(ref.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), ""))
	} else {
		for _, msg := range pathvalidation.IsValidPathSegmentName(ref.Name) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("name"), ref.Name, msg))
		}
	}

	return allErrs
}

var validMetricSourceTypes = sets.NewString(string(autoscaling.ObjectMetricSourceType), string(autoscaling.PodsMetricSourceType), string(autoscaling.ResourceMetricSourceType), string(autoscaling.ExternalMetricSourceType))
var validMetricSourceTypesList = validMetricSourceTypes.List()

func validateMetricSpec(spec autoscaling.MetricSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(string(spec.Type)) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("type"), "must specify a metric source type"))
	}

	if !validMetricSourceTypes.Has(string(spec.Type)) {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("type"), spec.Type, validMetricSourceTypesList))
	}

	typesPresent := sets.NewString()
	if spec.Object != nil {
		typesPresent.Insert("object")
		if typesPresent.Len() == 1 {
			allErrs = append(allErrs, validateObjectSource(spec.Object, fldPath.Child("object"))...)
		}
	}

	if spec.External != nil {
		typesPresent.Insert("external")
		if typesPresent.Len() == 1 {
			allErrs = append(allErrs, validateExternalSource(spec.External, fldPath.Child("external"))...)
		}
	}

	if spec.Pods != nil {
		typesPresent.Insert("pods")
		if typesPresent.Len() == 1 {
			allErrs = append(allErrs, validatePodsSource(spec.Pods, fldPath.Child("pods"))...)
		}
	}

	if spec.Resource != nil {
		typesPresent.Insert("resource")
		if typesPresent.Len() == 1 {
			allErrs = append(allErrs, validateResourceSource(spec.Resource, fldPath.Child("resource"))...)
		}
	}

	expectedField := strings.ToLower(string(spec.Type))

	if !typesPresent.Has(expectedField) {
		allErrs = append(allErrs, field.Required(fldPath.Child(expectedField), "must populate information for the given metric source"))
	}

	if typesPresent.Len() != 1 {
		typesPresent.Delete(expectedField)
		for typ := range typesPresent {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child(typ), "must populate the given metric source only"))
		}
	}

	return allErrs
}

func validateObjectSource(src *autoscaling.ObjectMetricSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, validateCrossVersionObjectReference(src.Target, fldPath.Child("target"))...)

	if len(src.MetricName) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("metricName"), "must specify a metric name"))
	}

	if src.TargetValue.Sign() != 1 {
		allErrs = append(allErrs, field.Required(fldPath.Child("targetValue"), "must specify a positive target value"))
	}

	return allErrs
}

func validateExternalSource(src *autoscaling.ExternalMetricSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(src.MetricName) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("metricName"), "must specify a metric name"))
	} else {
		for _, msg := range pathvalidation.IsValidPathSegmentName(src.MetricName) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("metricName"), src.MetricName, msg))
		}
	}

	if src.TargetValue == nil && src.TargetAverageValue == nil {
		allErrs = append(allErrs, field.Required(fldPath.Child("targetValue"), "must set either a target value for metric or a per-pod target"))
	}

	if src.TargetValue != nil && src.TargetAverageValue != nil {
		allErrs = append(allErrs, field.Forbidden(fldPath.Child("targetValue"), "may not set both a target value for metric and a per-pod target"))
	}

	if src.TargetAverageValue != nil && src.TargetAverageValue.Sign() != 1 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("targetAverageValue"), src.TargetAverageValue, "must be positive"))
	}

	if src.TargetValue != nil && src.TargetValue.Sign() != 1 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("targetValue"), src.TargetValue, "must be positive"))
	}

	return allErrs
}

func validatePodsSource(src *autoscaling.PodsMetricSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(src.MetricName) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("metricName"), "must specify a metric name"))
	}

	if src.TargetAverageValue.Sign() != 1 {
		allErrs = append(allErrs, field.Required(fldPath.Child("targetAverageValue"), "must specify a positive target value"))
	}

	return allErrs
}

func validateResourceSource(src *autoscaling.ResourceMetricSource, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(src.Name) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("name"), "must specify a resource name"))
	}

	if src.TargetAverageUtilization == nil && src.TargetAverageValue == nil {
		allErrs = append(allErrs, field.Required(fldPath.Child("targetAverageUtilization"), "must set either a target raw value or a target utilization"))
	}

	if src.TargetAverageUtilization != nil && *src.TargetAverageUtilization < 1 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("targetAverageUtilization"), src.TargetAverageUtilization, "must be greater than 0"))
	}

	if src.TargetAverageUtilization != nil && src.TargetAverageValue != nil {
		allErrs = append(allErrs, field.Forbidden(fldPath.Child("targetAverageValue"), "may not set both a target raw value and a target utilization"))
	}

	if src.TargetAverageValue != nil && src.TargetAverageValue.Sign() != 1 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("targetAverageValue"), src.TargetAverageValue, "must be positive"))
	}

	return allErrs
}
