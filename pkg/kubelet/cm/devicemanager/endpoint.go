package devicemanager

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"k8s.io/klog"

	pluginapi "k8s.io/kubelet/pkg/apis/deviceplugin/v1beta1"
)

// 	@implementBy: endpointImpl{}
//
// endpoint maps to a single registered device plugin. It is responsible
// for managing gRPC communications with the device plugin and caching
// device states reported by the device plugin.
type endpoint interface {
	run()
	stop()
	allocate(devs []string) (*pluginapi.AllocateResponse, error)
	preStartContainer(devs []string) (*pluginapi.PreStartContainerResponse, error)
	callback(resourceName string, devices []pluginapi.Device)
	isStopped() bool
	stopGracePeriodExpired() bool
}

type endpointImpl struct {
	client     pluginapi.DevicePluginClient
	clientConn *grpc.ClientConn

	socketPath   string
	resourceName string
	stopTime     time.Time

	mutex sync.Mutex
	cb    monitorCallback
}

// 	@param socketPath: device plugin 的 socket 路径, 在 /var/lib/kubelet/device-plugins 目录下
// 	@param resourceName: device plugin 提供的扩展资源名称
// 	@param callback: pkg/kubelet/cm/devicemanager/manager.go -> ManagerImpl.genericDeviceUpdateCallback()
//
// caller:
// 	1. pkg/kubelet/cm/devicemanager/manager.go -> ManagerImpl.addEndpoint()
//  kubelet 接收到来自 device plugin 的注册请求时被调用
//
// newEndpointImpl creates a new endpoint for the given resourceName.
// This is to be used during normal device plugin registration.
func newEndpointImpl(
	socketPath, resourceName string, callback monitorCallback,
) (*endpointImpl, error) {
	client, c, err := dial(socketPath)
	if err != nil {
		klog.Errorf("Can't create new endpoint with path %s err %v", socketPath, err)
		return nil, err
	}

	return &endpointImpl{
		client:     client,
		clientConn: c,

		socketPath:   socketPath,
		resourceName: resourceName,

		cb: callback,
	}, nil
}

// newStoppedEndpointImpl creates a new endpoint for the given resourceName with stopTime set.
// This is to be used during Kubelet restart, before the actual device plugin re-registers.
func newStoppedEndpointImpl(resourceName string) *endpointImpl {
	return &endpointImpl{
		resourceName: resourceName,
		stopTime:     time.Now(),
	}
}

// callback kubelet 通过 ListAndWatch 从 device plugin 处获得本机上 devicce 设备
// 的变动后, 调用该方法进行回调更新.
func (e *endpointImpl) callback(resourceName string, devices []pluginapi.Device) {
	e.cb(resourceName, devices)
}

// caller:
// 	1. pkg/kubelet/cm/devicemanager/manager.go -> ManagerImpl.runEndpoint()
//
// run initializes ListAndWatch gRPC call for the device plugin and
// blocks on receiving ListAndWatch gRPC stream updates.
// Each ListAndWatch stream update contains a new list of device states.
// It then issues a callback to pass this information to the device manager which
// will adjust the resource available information accordingly.
func (e *endpointImpl) run() {
	stream, err := e.client.ListAndWatch(context.Background(), &pluginapi.Empty{})
	if err != nil {
		klog.Errorf(errListAndWatch, e.resourceName, err)

		return
	}

	for {
		response, err := stream.Recv()
		if err != nil {
			klog.Errorf(errListAndWatch, e.resourceName, err)
			return
		}

		devs := response.Devices
		klog.V(2).Infof("State pushed for device plugin %s", e.resourceName)

		var newDevs []pluginapi.Device
		for _, d := range devs {
			newDevs = append(newDevs, *d)
		}

		e.callback(e.resourceName, newDevs)
	}
}

func (e *endpointImpl) isStopped() bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return !e.stopTime.IsZero()
}

func (e *endpointImpl) stopGracePeriodExpired() bool {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return !e.stopTime.IsZero() && time.Since(e.stopTime) > endpointStopGracePeriod
}

// used for testing only
func (e *endpointImpl) setStopTime(t time.Time) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.stopTime = t
}

// allocate 在 kubelet 为某个 Pod 分配好扩展资源的 deviceID 列表后, 需要额外调用
// device plugin 获取这些 device 的详细规格, 用于创建 runc 容器时指定 device 设备列表,
// 以及 mount 挂载列表.
//
// 	@param devs: 将要分配给目标 pod 的设备列表.
//
// caller:
// 	1. pkg/kubelet/cm/devicemanager/manager.go -> ManagerImpl.allocateContainerResources()
//
// allocate issues Allocate gRPC call to the device plugin.
func (e *endpointImpl) allocate(devs []string) (*pluginapi.AllocateResponse, error) {
	if e.isStopped() {
		return nil, fmt.Errorf(errEndpointStopped, e)
	}
	return e.client.Allocate(context.Background(), &pluginapi.AllocateRequest{
		ContainerRequests: []*pluginapi.ContainerAllocateRequest{
			{DevicesIDs: devs},
		},
	})
}

// preStartContainer issues PreStartContainer gRPC call to the device plugin.
func (e *endpointImpl) preStartContainer(devs []string) (*pluginapi.PreStartContainerResponse, error) {
	if e.isStopped() {
		return nil, fmt.Errorf(errEndpointStopped, e)
	}
	ctx, cancel := context.WithTimeout(
		context.Background(), pluginapi.KubeletPreStartContainerRPCTimeoutInSecs*time.Second,
	)
	defer cancel()
	return e.client.PreStartContainer(ctx, &pluginapi.PreStartContainerRequest{
		DevicesIDs: devs,
	})
}

func (e *endpointImpl) stop() {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.clientConn != nil {
		e.clientConn.Close()
	}
	e.stopTime = time.Now()
}

// dial establishes the gRPC communication with the registered device plugin. https://godoc.org/google.golang.org/grpc#Dial
func dial(unixSocketPath string) (pluginapi.DevicePluginClient, *grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c, err := grpc.DialContext(ctx, unixSocketPath, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithDialer(func(addr string, timeout time.Duration) (net.Conn, error) {
			return net.DialTimeout("unix", addr, timeout)
		}),
	)

	if err != nil {
		return nil, nil, fmt.Errorf(errFailedToDialDevicePlugin+" %v", err)
	}

	return pluginapi.NewDevicePluginClient(c), c, nil
}
