package grpcproxy

import (
	"context"
	"net/http"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/etcdhttp"
	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"
)

// HandleHealth registers health handler on '/health'.
func HandleHealth(mux *http.ServeMux, c *clientv3.Client) {
	mux.Handle(etcdhttp.PathHealth, etcdhttp.NewHealthHandler(func() etcdhttp.Health { return checkHealth(c) }))
}

func checkHealth(c *clientv3.Client) etcdhttp.Health {
	h := etcdhttp.Health{Health: "false"}
	ctx, cancel := context.WithTimeout(c.Ctx(), time.Second)
	_, err := c.Get(ctx, "a")
	cancel()
	if err == nil || err == rpctypes.ErrPermissionDenied {
		h.Health = "true"
	}
	return h
}
