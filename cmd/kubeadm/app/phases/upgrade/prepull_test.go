package upgrade

import (
	"testing"
	"time"

	"github.com/pkg/errors"

	"k8s.io/kubernetes/cmd/kubeadm/app/constants"
	//"k8s.io/apimachinery/pkg/util/version"
)

// failedCreatePrepuller is a fake prepuller that errors for kube-controller-manager in the CreateFunc call
type failedCreatePrepuller struct{}

func NewFailedCreatePrepuller() Prepuller {
	return &failedCreatePrepuller{}
}

func (p *failedCreatePrepuller) CreateFunc(component string) error {
	if component == "kube-controller-manager" {
		return errors.New("boo")
	}
	return nil
}

func (p *failedCreatePrepuller) WaitFunc(component string) {}

func (p *failedCreatePrepuller) DeleteFunc(component string) error {
	return nil
}

// foreverWaitPrepuller is a fake prepuller that basically waits "forever" (10 mins, but longer than the 10sec timeout)
type foreverWaitPrepuller struct{}

func NewForeverWaitPrepuller() Prepuller {
	return &foreverWaitPrepuller{}
}

func (p *foreverWaitPrepuller) CreateFunc(component string) error {
	return nil
}

func (p *foreverWaitPrepuller) WaitFunc(component string) {
	time.Sleep(10 * time.Minute)
}

func (p *foreverWaitPrepuller) DeleteFunc(component string) error {
	return nil
}

// failedDeletePrepuller is a fake prepuller that errors for kube-scheduler in the DeleteFunc call
type failedDeletePrepuller struct{}

func NewFailedDeletePrepuller() Prepuller {
	return &failedDeletePrepuller{}
}

func (p *failedDeletePrepuller) CreateFunc(component string) error {
	return nil
}

func (p *failedDeletePrepuller) WaitFunc(component string) {}

func (p *failedDeletePrepuller) DeleteFunc(component string) error {
	if component == "kube-scheduler" {
		return errors.New("boo")
	}
	return nil
}

// goodPrepuller is a fake prepuller that works as expected
type goodPrepuller struct{}

func NewGoodPrepuller() Prepuller {
	return &goodPrepuller{}
}

func (p *goodPrepuller) CreateFunc(component string) error {
	time.Sleep(300 * time.Millisecond)
	return nil
}

func (p *goodPrepuller) WaitFunc(component string) {
	time.Sleep(300 * time.Millisecond)
}

func (p *goodPrepuller) DeleteFunc(component string) error {
	time.Sleep(300 * time.Millisecond)
	return nil
}

func TestPrepullImagesInParallel(t *testing.T) {
	tests := []struct {
		name        string
		p           Prepuller
		timeout     time.Duration
		expectedErr bool
	}{
		{
			name:        "should error out; create failed",
			p:           NewFailedCreatePrepuller(),
			timeout:     constants.PrepullImagesInParallelTimeout,
			expectedErr: true,
		},
		{
			name:        "should error out; timeout exceeded",
			p:           NewForeverWaitPrepuller(),
			timeout:     constants.PrepullImagesInParallelTimeout,
			expectedErr: true,
		},
		{
			name:        "should error out; delete failed",
			p:           NewFailedDeletePrepuller(),
			timeout:     constants.PrepullImagesInParallelTimeout,
			expectedErr: true,
		},
		{
			name:        "should work just fine",
			p:           NewGoodPrepuller(),
			timeout:     constants.PrepullImagesInParallelTimeout,
			expectedErr: false,
		},
	}

	for _, rt := range tests {
		t.Run(rt.name, func(t *testing.T) {
			actualErr := PrepullImagesInParallel(rt.p, rt.timeout, append(constants.ControlPlaneComponents, constants.Etcd))
			if (actualErr != nil) != rt.expectedErr {
				t.Errorf(
					"failed TestPrepullImagesInParallel\n\texpected error: %t\n\tgot: %t",
					rt.expectedErr,
					(actualErr != nil),
				)
			}
		})
	}
}
