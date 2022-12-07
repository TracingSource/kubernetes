package fuzzer

import (
	fuzz "github.com/google/gofuzz"

	runtimeserializer "k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kubernetes/pkg/apis/auditregistration"
)

// Funcs returns the fuzzer functions for the auditregistration api group.
var Funcs = func(codecs runtimeserializer.CodecFactory) []interface{} {
	return []interface{}{
		func(obj *auditregistration.AuditSink, c fuzz.Continue) {
			c.FuzzNoCustom(obj)
			v := int64(1)
			obj.Spec.Webhook.Throttle = &auditregistration.WebhookThrottleConfig{
				QPS:   &v,
				Burst: &v,
			}
		},
	}
}
