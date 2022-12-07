package v1beta1

import (
	"fmt"

	appsv1beta1 "k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/apis/autoscaling"
)

func addConversionFuncs(scheme *runtime.Scheme) error {
	// Add field label conversions for kinds having selectable nothing but ObjectMeta fields.
	if err := scheme.AddFieldLabelConversionFunc(SchemeGroupVersion.WithKind("StatefulSet"),
		func(label, value string) (string, string, error) {
			switch label {
			case "metadata.name", "metadata.namespace", "status.successful":
				return label, value, nil
			default:
				return "", "", fmt.Errorf("field label not supported for appsv1beta1.StatefulSet: %s", label)
			}
		}); err != nil {
		return err
	}

	return nil
}

func Convert_autoscaling_ScaleStatus_To_v1beta1_ScaleStatus(in *autoscaling.ScaleStatus, out *appsv1beta1.ScaleStatus, s conversion.Scope) error {
	out.Replicas = int32(in.Replicas)
	out.TargetSelector = in.Selector

	out.Selector = nil
	selector, err := metav1.ParseToLabelSelector(in.Selector)
	if err != nil {
		return fmt.Errorf("failed to parse selector: %v", err)
	}
	if len(selector.MatchExpressions) == 0 {
		out.Selector = selector.MatchLabels
	}

	return nil
}

func Convert_v1beta1_ScaleStatus_To_autoscaling_ScaleStatus(in *appsv1beta1.ScaleStatus, out *autoscaling.ScaleStatus, s conversion.Scope) error {
	out.Replicas = in.Replicas

	if in.TargetSelector != "" {
		out.Selector = in.TargetSelector
	} else if in.Selector != nil {
		set := labels.Set{}
		for key, val := range in.Selector {
			set[key] = val
		}
		out.Selector = labels.SelectorFromSet(set).String()
	} else {
		out.Selector = ""
	}
	return nil
}
