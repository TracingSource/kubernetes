package devicemapper

import (
	"os/exec"
	"strconv"
	"strings"

	"k8s.io/klog"
)

// DmsetupClient is a low-level client for interacting with device mapper via
// the `dmsetup` utility, which is provided by the `device-mapper` package.
type DmsetupClient interface {
	// Table runs `dmsetup table` on the given device name and returns the
	// output or an error.
	Table(deviceName string) ([]byte, error)
	// Message runs `dmsetup message` on the given device, passing the given
	// message to the given sector, and returns the output or an error.
	Message(deviceName string, sector int, message string) ([]byte, error)
	// Status runs `dmsetup status` on the given device and returns the output
	// or an error.
	Status(deviceName string) ([]byte, error)
}

// NewDmSetupClient returns a new DmsetupClient.
func NewDmsetupClient() DmsetupClient {
	return &defaultDmsetupClient{}
}

// defaultDmsetupClient is a functional DmsetupClient
type defaultDmsetupClient struct{}

var _ DmsetupClient = &defaultDmsetupClient{}

func (c *defaultDmsetupClient) Table(deviceName string) ([]byte, error) {
	return c.dmsetup("table", deviceName)
}

func (c *defaultDmsetupClient) Message(deviceName string, sector int, message string) ([]byte, error) {
	return c.dmsetup("message", deviceName, strconv.Itoa(sector), message)
}

func (c *defaultDmsetupClient) Status(deviceName string) ([]byte, error) {
	return c.dmsetup("status", deviceName)
}

func (*defaultDmsetupClient) dmsetup(args ...string) ([]byte, error) {
	klog.V(5).Infof("running dmsetup %v", strings.Join(args, " "))
	return exec.Command("dmsetup", args...).Output()
}
