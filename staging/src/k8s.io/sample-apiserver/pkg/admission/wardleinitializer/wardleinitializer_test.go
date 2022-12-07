package wardleinitializer_test

import (
	"context"
	"testing"
	"time"

	"k8s.io/apiserver/pkg/admission"
	"k8s.io/sample-apiserver/pkg/admission/wardleinitializer"
	"k8s.io/sample-apiserver/pkg/generated/clientset/versioned/fake"
	informers "k8s.io/sample-apiserver/pkg/generated/informers/externalversions"
)

// TestWantsInternalWardleInformerFactory ensures that the informer factory is injected
// when the WantsInternalWardleInformerFactory interface is implemented by a plugin.
func TestWantsInternalWardleInformerFactory(t *testing.T) {
	cs := &fake.Clientset{}
	sf := informers.NewSharedInformerFactory(cs, time.Duration(1)*time.Second)
	target := wardleinitializer.New(sf)

	wantWardleInformerFactory := &wantInternalWardleInformerFactory{}
	target.Initialize(wantWardleInformerFactory)
	if wantWardleInformerFactory.sf != sf {
		t.Errorf("expected informer factory to be initialized")
	}
}

// wantInternalWardleInformerFactory is a test stub that fulfills the WantsInternalWardleInformerFactory interface
type wantInternalWardleInformerFactory struct {
	sf informers.SharedInformerFactory
}

func (self *wantInternalWardleInformerFactory) SetInternalWardleInformerFactory(sf informers.SharedInformerFactory) {
	self.sf = sf
}
func (self *wantInternalWardleInformerFactory) Admit(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	return nil
}
func (self *wantInternalWardleInformerFactory) Handles(o admission.Operation) bool { return false }
func (self *wantInternalWardleInformerFactory) ValidateInitialization() error      { return nil }

var _ admission.Interface = &wantInternalWardleInformerFactory{}
var _ wardleinitializer.WantsInternalWardleInformerFactory = &wantInternalWardleInformerFactory{}
