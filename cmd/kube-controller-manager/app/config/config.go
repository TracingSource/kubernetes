package config

import (
	apiserver "k8s.io/apiserver/pkg/server"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
	kubectrlmgrconfig "k8s.io/kubernetes/pkg/controller/apis/config"
)

// Config is the main context object for the controller manager.
type Config struct {
	// ComponentConfig controller manager程序所需的参数, config还有其他的配置.
	ComponentConfig kubectrlmgrconfig.KubeControllerManagerConfiguration

	SecureServing *apiserver.SecureServingInfo
	// LoopbackClientConfig is a config for a privileged loopback connection
	LoopbackClientConfig *restclient.Config

	// TODO: remove deprecated insecure serving
	InsecureServing *apiserver.DeprecatedInsecureServingInfo
	// Authentication http 服务的认证中间件
	// 在 cmd/kube-controller-manager/app/controllermanager.go -> Run() 中使用(BuildHandlerChain)
	//
	// Authentication 通过如下文件完成初始化
	// staging/src/k8s.io/apiserver/pkg/server/options/authentication.go ->
	// DelegatingAuthenticationOptions.ApplyTo()
	Authentication  apiserver.AuthenticationInfo
	// Authorization http 服务的认证中间件
	// 在 cmd/kube-controller-manager/app/controllermanager.go -> Run() 中使用(BuildHandlerChain)
	//
	// Authorization 通过如下文件完成初始化
	// staging/src/k8s.io/apiserver/pkg/server/options/authorization.go ->
	// DelegatingAuthorizationOptions.ApplyTo()
	Authorization   apiserver.AuthorizationInfo

	// Client 在如下文件中被初始化并赋值
	// cmd/kube-controller-manager/app/options/options.go -> 
	// KubeControllerManagerOptions.Config()
	//
	// the general kube client
	Client *clientset.Clientset

	// the client only used for leader election
	LeaderElectionClient *clientset.Clientset

	// the rest config for the master
	Kubeconfig *restclient.Config

	// the event sink
	EventRecorder record.EventRecorder
}

type completedConfig struct {
	*Config
}

// 调用者无法直接使用其成员, 只能调用Complete()方法.
//
// CompletedConfig same as Config, just to swap private object.
type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// caller: 
// 	1. cmd/kube-controller-manager/app/config/config.go -> NewControllerManagerCommand()
//
// Complete fills in any fields not set that are required to have valid data.
// It's mutating the receiver.
func (c *Config) Complete() *CompletedConfig {
	cc := completedConfig{c}

	apiserver.AuthorizeClientBearerToken(
		c.LoopbackClientConfig, &c.Authentication, &c.Authorization,
	)

	return &CompletedConfig{&cc}
}
