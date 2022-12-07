package grpcproxy

import (
	"context"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/etcdserver/api/v3lock/v3lockpb"
)

type lockProxy struct {
	client *clientv3.Client
}

func NewLockProxy(client *clientv3.Client) v3lockpb.LockServer {
	return &lockProxy{client: client}
}

func (lp *lockProxy) Lock(ctx context.Context, req *v3lockpb.LockRequest) (*v3lockpb.LockResponse, error) {
	return v3lockpb.NewLockClient(lp.client.ActiveConnection()).Lock(ctx, req)
}

func (lp *lockProxy) Unlock(ctx context.Context, req *v3lockpb.UnlockRequest) (*v3lockpb.UnlockResponse, error) {
	return v3lockpb.NewLockClient(lp.client.ActiveConnection()).Unlock(ctx, req)
}
