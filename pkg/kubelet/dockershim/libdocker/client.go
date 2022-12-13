package libdocker

import (
	"time"

	dockertypes "github.com/docker/docker/api/types"
	dockercontainer "github.com/docker/docker/api/types/container"
	dockerimagetypes "github.com/docker/docker/api/types/image"
	dockerapi "github.com/docker/docker/client"
	"k8s.io/klog"
)

const (
	// https://docs.docker.com/engine/reference/api/docker_remote_api/
	// docker version should be at least 1.13.1
	MinimumDockerAPIVersion = "1.26.0"

	// Status of a container returned by ListContainers.
	StatusRunningPrefix = "Up"
	StatusCreatedPrefix = "Created"
	StatusExitedPrefix  = "Exited"

	// Fake docker endpoint
	FakeDockerEndpoint = "fake://"
)

// Interface is an abstract interface for testability. 
// It abstracts the interface of docker client.
type Interface interface {
	ListContainers(options dockertypes.ContainerListOptions) ([]dockertypes.Container, error)
	InspectContainer(id string) (*dockertypes.ContainerJSON, error)
	InspectContainerWithSize(id string) (*dockertypes.ContainerJSON, error)
	CreateContainer(dockertypes.ContainerCreateConfig) (*dockercontainer.ContainerCreateCreatedBody, error)
	StartContainer(id string) error
	StopContainer(id string, timeout time.Duration) error
	UpdateContainerResources(id string, updateConfig dockercontainer.UpdateConfig) error
	RemoveContainer(id string, opts dockertypes.ContainerRemoveOptions) error
	InspectImageByRef(imageRef string) (*dockertypes.ImageInspect, error)
	InspectImageByID(imageID string) (*dockertypes.ImageInspect, error)
	ListImages(opts dockertypes.ImageListOptions) ([]dockertypes.ImageSummary, error)
	PullImage(image string, auth dockertypes.AuthConfig, opts dockertypes.ImagePullOptions) error
	RemoveImage(image string, opts dockertypes.ImageRemoveOptions) ([]dockertypes.ImageDeleteResponseItem, error)
	ImageHistory(id string) ([]dockerimagetypes.HistoryResponseItem, error)
	Logs(string, dockertypes.ContainerLogsOptions, StreamOptions) error
	Version() (*dockertypes.Version, error)
	Info() (*dockertypes.Info, error)
	CreateExec(string, dockertypes.ExecConfig) (*dockertypes.IDResponse, error)
	StartExec(string, dockertypes.ExecStartCheck, StreamOptions) error
	InspectExec(id string) (*dockertypes.ContainerExecInspect, error)
	AttachToContainer(string, dockertypes.ContainerAttachOptions, StreamOptions) error
	ResizeContainerTTY(id string, height, width uint) error
	ResizeExecTTY(id string, height, width uint) error
	GetContainerStats(id string) (*dockertypes.StatsJSON, error)
}

// getDockerClient 调用 docker 官方库, 创建与 dockerd 服务进行通信的客户端并返回
//
// 	@param dockerEndpoint: unix:///var/run/docker.sock
//
// Get a *dockerapi.Client, either using the endpoint passed in,
// or using DOCKER_HOST, DOCKER_TLS_VERIFY, and DOCKER_CERT path per their spec
func getDockerClient(dockerEndpoint string) (*dockerapi.Client, error) {
	if len(dockerEndpoint) > 0 {
		klog.Infof("Connecting to docker on %s", dockerEndpoint)
		return dockerapi.NewClient(dockerEndpoint, "", nil, nil)
	}
	return dockerapi.NewClientWithOpts(dockerapi.FromEnv)
}

// ConnectToDockerOrDie 调用 docker 官方库, 创建与 dockerd 服务进行通信的客户端并返回
//
// 	@param dockerEndpoint: unix:///var/run/docker.sock
//
// caller:
// 	1. pkg/kubelet/dockershim/docker_service.go -> NewDockerClientFromConfig()
// 	kubelet 启动时, 建立与 dockerd 服务的通信, 同时启动 dockershim 进程, 
// 	将通用的 cri 请求进行转换, 转换成 dockerd 服务接受的形式.
//
// ConnectToDockerOrDie creates docker client connecting to docker daemon.
// If the endpoint passed in is "fake://", a fake docker client will be returned.
// The program exits if error occurs.
// The requestTimeout is the timeout for docker requests.
// If timeout is exceeded, the request will be cancelled and throw out an error.
// If requestTimeout is 0, a default value will be applied.
func ConnectToDockerOrDie(
	dockerEndpoint string, 
	requestTimeout, imagePullProgressDeadline time.Duration,
	withTraceDisabled bool, enableSleep bool,
) Interface {
	if dockerEndpoint == FakeDockerEndpoint {
		fakeClient := NewFakeDockerClient()
		if withTraceDisabled {
			fakeClient = fakeClient.WithTraceDisabled()
		}

		if enableSleep {
			fakeClient.EnableSleep = true
		}
		return fakeClient
	}
	client, err := getDockerClient(dockerEndpoint)
	if err != nil {
		klog.Fatalf("Couldn't connect to docker: %v", err)
	}
	klog.Infof("Start docker client with request timeout=%v", requestTimeout)
	return newKubeDockerClient(client, requestTimeout, imagePullProgressDeadline)
}
