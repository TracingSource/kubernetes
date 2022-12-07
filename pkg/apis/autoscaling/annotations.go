package autoscaling

// MetricSpecsAnnotation is the annotation which holds non-CPU-utilization HPA metric
// specs when converting the `Metrics` field from autoscaling/v2beta1
const MetricSpecsAnnotation = "autoscaling.alpha.kubernetes.io/metrics"

// MetricStatusesAnnotation is the annotation which holds non-CPU-utilization HPA metric
// statuses when converting the `CurrentMetrics` field from autoscaling/v2beta1
const MetricStatusesAnnotation = "autoscaling.alpha.kubernetes.io/current-metrics"

// HorizontalPodAutoscalerConditionsAnnotation is the annotation which holds the conditions
// of an HPA when converting the `Conditions` field from autoscaling/v2beta1
const HorizontalPodAutoscalerConditionsAnnotation = "autoscaling.alpha.kubernetes.io/conditions"

// DefaultCPUUtilization is the default value for CPU utilization, provided no other
// metrics are present.  This is here because it's used by both the v2beta1 defaulting
// logic, and the pseudo-defaulting done in v1 conversion.
const DefaultCPUUtilization = 80
