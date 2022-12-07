package checkpoint

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
)

// configMapPayload implements Payload, backed by a v1/ConfigMap config source object
type configMapPayload struct {
	cm *apiv1.ConfigMap
}

var _ Payload = (*configMapPayload)(nil)

// NewConfigMapPayload constructs a Payload backed by a ConfigMap, which must have a non-empty UID
func NewConfigMapPayload(cm *apiv1.ConfigMap) (Payload, error) {
	if cm == nil {
		return nil, fmt.Errorf("ConfigMap must be non-nil")
	} else if cm.ObjectMeta.UID == "" {
		return nil, fmt.Errorf("ConfigMap must have a non-empty UID")
	} else if cm.ObjectMeta.ResourceVersion == "" {
		return nil, fmt.Errorf("ConfigMap must have a non-empty ResourceVersion")
	}

	return &configMapPayload{cm}, nil
}

func (p *configMapPayload) UID() string {
	return string(p.cm.UID)
}

func (p *configMapPayload) ResourceVersion() string {
	return p.cm.ResourceVersion
}

func (p *configMapPayload) Files() map[string]string {
	return p.cm.Data
}

func (p *configMapPayload) object() interface{} {
	return p.cm
}
