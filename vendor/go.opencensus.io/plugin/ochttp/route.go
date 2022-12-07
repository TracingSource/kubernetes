package ochttp

import (
	"context"
	"net/http"

	"go.opencensus.io/tag"
)

// SetRoute sets the http_server_route tag to the given value.
// It's useful when an HTTP framework does not support the http.Handler interface
// and using WithRouteTag is not an option, but provides a way to hook into the request flow.
func SetRoute(ctx context.Context, route string) {
	if a, ok := ctx.Value(addedTagsKey{}).(*addedTags); ok {
		a.t = append(a.t, tag.Upsert(KeyServerRoute, route))
	}
}

// WithRouteTag returns an http.Handler that records stats with the
// http_server_route tag set to the given value.
func WithRouteTag(handler http.Handler, route string) http.Handler {
	return taggedHandlerFunc(func(w http.ResponseWriter, r *http.Request) []tag.Mutator {
		addRoute := []tag.Mutator{tag.Upsert(KeyServerRoute, route)}
		ctx, _ := tag.New(r.Context(), addRoute...)
		r = r.WithContext(ctx)
		handler.ServeHTTP(w, r)
		return addRoute
	})
}

// taggedHandlerFunc is a http.Handler that returns tags describing the
// processing of the request. These tags will be recorded along with the
// measures in this package at the end of the request.
type taggedHandlerFunc func(w http.ResponseWriter, r *http.Request) []tag.Mutator

func (h taggedHandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tags := h(w, r)
	if a, ok := r.Context().Value(addedTagsKey{}).(*addedTags); ok {
		a.t = append(a.t, tags...)
	}
}

type addedTagsKey struct{}

type addedTags struct {
	t []tag.Mutator
}
