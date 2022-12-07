package v3client

import (
	"context"
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver"
	"go.etcd.io/etcd/etcdserver/api/v3rpc"
	"go.etcd.io/etcd/proxy/grpcproxy/adapter"
)

// New creates a clientv3 client that wraps an in-process EtcdServer. Instead
// of making gRPC calls through sockets, the client makes direct function calls
// to the etcd server through its api/v3rpc function interfaces.
func New(s *etcdserver.EtcdServer) *clientv3.Client {
	c := clientv3.NewCtxClient(context.Background())

	kvc := adapter.KvServerToKvClient(v3rpc.NewQuotaKVServer(s))
	c.KV = clientv3.NewKVFromKVClient(kvc, c)

	lc := adapter.LeaseServerToLeaseClient(v3rpc.NewQuotaLeaseServer(s))
	c.Lease = clientv3.NewLeaseFromLeaseClient(lc, c, time.Second)

	wc := adapter.WatchServerToWatchClient(v3rpc.NewWatchServer(s))
	c.Watcher = &watchWrapper{clientv3.NewWatchFromWatchClient(wc, c)}

	mc := adapter.MaintenanceServerToMaintenanceClient(v3rpc.NewMaintenanceServer(s))
	c.Maintenance = clientv3.NewMaintenanceFromMaintenanceClient(mc, c)

	clc := adapter.ClusterServerToClusterClient(v3rpc.NewClusterServer(s))
	c.Cluster = clientv3.NewClusterFromClusterClient(clc, c)

	// TODO: implement clientv3.Auth interface?

	return c
}

// BlankContext implements Stringer on a context so the ctx string doesn't
// depend on the context's WithValue data, which tends to be unsynchronized
// (e.g., x/net/trace), causing ctx.String() to throw data races.
type blankContext struct{ context.Context }

func (*blankContext) String() string { return "(blankCtx)" }

// watchWrapper wraps clientv3 watch calls to blank out the context
// to avoid races on trace data.
type watchWrapper struct{ clientv3.Watcher }

func (ww *watchWrapper) Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	return ww.Watcher.Watch(&blankContext{ctx}, key, opts...)
}
