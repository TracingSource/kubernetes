// Package v3client provides clientv3 interfaces from an etcdserver.
//
// Use v3client by creating an EtcdServer instance, then wrapping it with v3client.New:
//
//	import (
//		"context"
//
//		"go.etcd.io/etcd/embed"
//		"go.etcd.io/etcd/etcdserver/api/v3client"
//	)
//
//	...
//
//	// create an embedded EtcdServer from the default configuration
//	cfg := embed.NewConfig()
//	cfg.Dir = "default.etcd"
//	e, err := embed.StartEtcd(cfg)
//	if err != nil {
//		// handle error!
//	}
//
//	// wrap the EtcdServer with v3client
//	cli := v3client.New(e.Server)
//
//	// use like an ordinary clientv3
//	resp, err := cli.Put(context.TODO(), "some-key", "it works!")
//	if err != nil {
//		// handle error!
//	}
//
package v3client
