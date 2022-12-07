package admission

import (
	"context"
	"testing"

	"k8s.io/apiserver/pkg/admission"
)

type doNothingAdmission struct{}

func (doNothingAdmission) Admit(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	return nil
}
func (doNothingAdmission) Handles(o admission.Operation) bool { return false }
func (doNothingAdmission) Validate() error                    { return nil }

type WantsCloudConfigAdmissionPlugin struct {
	doNothingAdmission
	cloudConfig []byte
}

func (p *WantsCloudConfigAdmissionPlugin) SetCloudConfig(cloudConfig []byte) {
	p.cloudConfig = cloudConfig
}

func TestCloudConfigAdmissionPlugin(t *testing.T) {
	cloudConfig := []byte("cloud-configuration")
	initializer := NewPluginInitializer(cloudConfig, nil, nil)
	wantsCloudConfigAdmission := &WantsCloudConfigAdmissionPlugin{}
	initializer.Initialize(wantsCloudConfigAdmission)

	if wantsCloudConfigAdmission.cloudConfig == nil {
		t.Errorf("Expected cloud config to be initialized but found nil")
	}
}
