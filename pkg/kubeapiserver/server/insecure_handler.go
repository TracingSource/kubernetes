package server

import (
	"net/http"

	genericapifilters "k8s.io/apiserver/pkg/endpoints/filters"
	"k8s.io/apiserver/pkg/server"
	genericfilters "k8s.io/apiserver/pkg/server/filters"
)

// BuildInsecureHandlerChain sets up the server to listen to http. Should be removed.
func BuildInsecureHandlerChain(apiHandler http.Handler, c *server.Config) http.Handler {
	handler := apiHandler
	handler = genericfilters.WithMaxInFlightLimit(handler, c.MaxRequestsInFlight, c.MaxMutatingRequestsInFlight, c.LongRunningFunc)
	handler = genericapifilters.WithAudit(handler, c.AuditBackend, c.AuditPolicyChecker, c.LongRunningFunc)
	handler = genericapifilters.WithAuthentication(handler, server.InsecureSuperuser{}, nil, nil)
	handler = genericfilters.WithCORS(handler, c.CorsAllowedOriginList, nil, nil, nil, "true")
	handler = genericfilters.WithTimeoutForNonLongRunningRequests(handler, c.LongRunningFunc, c.RequestTimeout)
	handler = genericfilters.WithWaitGroup(handler, c.LongRunningFunc, c.HandlerChainWaitGroup)
	handler = genericapifilters.WithRequestInfo(handler, server.NewRequestInfoResolver(c))
	handler = genericapifilters.WithCacheControl(handler)
	handler = genericfilters.WithPanicRecovery(handler)

	return handler
}
