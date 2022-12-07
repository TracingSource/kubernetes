// Copyright 2016 Google LLC
//


// +build appengine

package http

import (
	"context"
	"net/http"

	"google.golang.org/appengine/urlfetch"
)

func init() {
	appengineUrlfetchHook = func(ctx context.Context) http.RoundTripper {
		return &urlfetch.Transport{Context: ctx}
	}
}
