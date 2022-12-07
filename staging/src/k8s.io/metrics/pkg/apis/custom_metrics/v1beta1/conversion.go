package v1beta1

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/conversion"
	"k8s.io/metrics/pkg/apis/custom_metrics"
)

func Convert_v1beta1_MetricValue_To_custom_metrics_MetricValue(in *MetricValue, out *custom_metrics.MetricValue, s conversion.Scope) error {
	out.TypeMeta = in.TypeMeta
	out.DescribedObject = custom_metrics.ObjectReference{
		Kind:            in.DescribedObject.Kind,
		Namespace:       in.DescribedObject.Namespace,
		Name:            in.DescribedObject.Name,
		UID:             in.DescribedObject.UID,
		APIVersion:      in.DescribedObject.APIVersion,
		ResourceVersion: in.DescribedObject.ResourceVersion,
		FieldPath:       in.DescribedObject.FieldPath,
	}
	out.Timestamp = in.Timestamp
	out.WindowSeconds = in.WindowSeconds
	out.Value = in.Value
	out.Metric.Name = in.MetricName
	out.Metric.Selector = in.Selector
	return nil
}

func Convert_custom_metrics_MetricValue_To_v1beta1_MetricValue(in *custom_metrics.MetricValue, out *MetricValue, s conversion.Scope) error {
	out.TypeMeta = in.TypeMeta
	out.DescribedObject = v1.ObjectReference{
		Kind:            in.DescribedObject.Kind,
		Namespace:       in.DescribedObject.Namespace,
		Name:            in.DescribedObject.Name,
		UID:             in.DescribedObject.UID,
		APIVersion:      in.DescribedObject.APIVersion,
		ResourceVersion: in.DescribedObject.ResourceVersion,
		FieldPath:       in.DescribedObject.FieldPath,
	}
	out.Timestamp = in.Timestamp
	out.WindowSeconds = in.WindowSeconds
	out.Value = in.Value
	out.MetricName = in.Metric.Name
	out.Selector = in.Metric.Selector
	return nil
}
