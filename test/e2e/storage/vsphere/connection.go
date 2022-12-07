package vsphere

import (
	"context"
	"fmt"
	neturl "net/url"
	"sync"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/session"
	"github.com/vmware/govmomi/vim25"
	"k8s.io/klog"
)

const (
	roundTripperDefaultCount = 3
)

var (
	clientLock sync.Mutex
)

// Connect makes connection to vSphere
// No actions are taken if a connection exists and alive. Otherwise, a new client will be created.
func Connect(ctx context.Context, vs *VSphere) error {
	var err error
	clientLock.Lock()
	defer clientLock.Unlock()

	if vs.Client == nil {
		vs.Client, err = NewClient(ctx, vs)
		if err != nil {
			klog.Errorf("Failed to create govmomi client. err: %+v", err)
			return err
		}
		return nil
	}
	manager := session.NewManager(vs.Client.Client)
	userSession, err := manager.UserSession(ctx)
	if err != nil {
		klog.Errorf("Error while obtaining user session. err: %+v", err)
		return err
	}
	if userSession != nil {
		return nil
	}
	klog.Warningf("Creating new client session since the existing session is not valid or not authenticated")
	vs.Client.Logout(ctx)
	vs.Client, err = NewClient(ctx, vs)
	if err != nil {
		klog.Errorf("Failed to create govmomi client. err: %+v", err)
		return err
	}
	return nil
}

// NewClient creates a new client for vSphere connection
func NewClient(ctx context.Context, vs *VSphere) (*govmomi.Client, error) {
	url, err := neturl.Parse(fmt.Sprintf("https://%s:%s/sdk", vs.Config.Hostname, vs.Config.Port))
	if err != nil {
		klog.Errorf("Failed to parse URL: %s. err: %+v", url, err)
		return nil, err
	}
	url.User = neturl.UserPassword(vs.Config.Username, vs.Config.Password)
	client, err := govmomi.NewClient(ctx, url, true)
	if err != nil {
		klog.Errorf("Failed to create new client. err: %+v", err)
		return nil, err
	}
	if vs.Config.RoundTripperCount == 0 {
		vs.Config.RoundTripperCount = roundTripperDefaultCount
	}
	client.RoundTripper = vim25.Retry(client.RoundTripper, vim25.TemporaryNetworkError(int(vs.Config.RoundTripperCount)))
	return client, nil
}
