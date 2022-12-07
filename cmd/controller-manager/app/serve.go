package app

import (
	"net/http"
	goruntime "runtime"

	genericapifilters "k8s.io/apiserver/pkg/endpoints/filters"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	apiserver "k8s.io/apiserver/pkg/server"
	genericfilters "k8s.io/apiserver/pkg/server/filters"
	"k8s.io/apiserver/pkg/server/healthz"
	"k8s.io/apiserver/pkg/server/mux"
	"k8s.io/apiserver/pkg/server/routes"
	componentbaseconfig "k8s.io/component-base/config"
	"k8s.io/component-base/metrics/legacyregistry"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/util/configz"
)

// BuildHandlerChain 为 http 服务器添加内置中间件(过滤器), 包含认证, 鉴权等.
//
// caller: 
// 	1. cmd/kube-controller-manager/app/controllermanager.go -> Run()
//
// BuildHandlerChain builds a handler chain with a base handler and CompletedConfig.
func BuildHandlerChain(
	apiHandler http.Handler, 
	authorizationInfo *apiserver.AuthorizationInfo, 
	authenticationInfo *apiserver.AuthenticationInfo,
) http.Handler {
	requestInfoResolver := &apirequest.RequestInfoFactory{}
	failedHandler := genericapifilters.Unauthorized(legacyscheme.Codecs, false)

	handler := apiHandler
	if authorizationInfo != nil {
		handler = genericapifilters.WithAuthorization(
			apiHandler, authorizationInfo.Authorizer, legacyscheme.Codecs,
		)
	}
	if authenticationInfo != nil {
		handler = genericapifilters.WithAuthentication(
			handler, authenticationInfo.Authenticator, failedHandler, nil,
		)
	}
	handler = genericapifilters.WithRequestInfo(handler, requestInfoResolver)
	handler = genericapifilters.WithCacheControl(handler)
	handler = genericfilters.WithPanicRecovery(handler)

	return handler
}

// NewBaseHandler 挂载 /healthz, /metrics, /configz 等接口, 并返回 mux 对象.
//
// caller: 
// 	1. cmd/kube-controller-manager/app/controllermanager.go -> Run()
//
// NewBaseHandler takes in CompletedConfig and returns a handler.
func NewBaseHandler(
	c *componentbaseconfig.DebuggingConfiguration, checks ...healthz.HealthChecker,
) *mux.PathRecorderMux {
	mux := mux.NewPathRecorderMux("controller-manager")
	healthz.InstallHandler(mux, checks...) // /healthz, /healthz/ping
	if c.EnableProfiling {
		routes.Profiling{}.Install(mux)
		if c.EnableContentionProfiling {
			goruntime.SetBlockProfileRate(1)
		}
	}
	configz.InstallHandler(mux) // /configz 
	//lint:ignore SA1019 See the Metrics Stability Migration KEP
	mux.Handle("/metrics", legacyregistry.Handler()) // /metrics

	return mux
}
