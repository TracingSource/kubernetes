// Package naming provides an etcd-backed gRPC resolver for discovering gRPC services.
//
// To use, first import the packages:
//
//	import (
//		"go.etcd.io/etcd/clientv3"
//		etcdnaming "go.etcd.io/etcd/clientv3/naming"
//
//		"google.golang.org/grpc"
//		"google.golang.org/grpc/naming"
//	)
//
// First, register new endpoint addresses for a service:
//
//	func etcdAdd(c *clientv3.Client, service, addr string) error {
//		r := &etcdnaming.GRPCResolver{Client: c}
//		return r.Update(c.Ctx(), service, naming.Update{Op: naming.Add, Addr: addr})
//	}
//
// Dial an RPC service using the etcd gRPC resolver and a gRPC Balancer:
//
//	func etcdDial(c *clientv3.Client, service string) (*grpc.ClientConn, error) {
//		r := &etcdnaming.GRPCResolver{Client: c}
//		b := grpc.RoundRobin(r)
//		return grpc.Dial(service, grpc.WithBalancer(b))
//	}
//
// Optionally, force delete an endpoint:
//
//	func etcdDelete(c *clientv3, service, addr string) error {
//		r := &etcdnaming.GRPCResolver{Client: c}
//		return r.Update(c.Ctx(), service, naming.Update{Op: naming.Delete, Addr: "1.2.3.4"})
//	}
//
// Or register an expiring endpoint with a lease:
//
//	func etcdLeaseAdd(c *clientv3.Client, lid clientv3.LeaseID, service, addr string) error {
//		r := &etcdnaming.GRPCResolver{Client: c}
//		return r.Update(c.Ctx(), service, naming.Update{Op: naming.Add, Addr: addr}, clientv3.WithLease(lid))
//	}
//
package naming
