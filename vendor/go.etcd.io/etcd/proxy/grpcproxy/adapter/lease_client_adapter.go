package adapter

import (
	"context"

	pb "go.etcd.io/etcd/etcdserver/etcdserverpb"

	"google.golang.org/grpc"
)

type ls2lc struct {
	leaseServer pb.LeaseServer
}

func LeaseServerToLeaseClient(ls pb.LeaseServer) pb.LeaseClient {
	return &ls2lc{ls}
}

func (c *ls2lc) LeaseGrant(ctx context.Context, in *pb.LeaseGrantRequest, opts ...grpc.CallOption) (*pb.LeaseGrantResponse, error) {
	return c.leaseServer.LeaseGrant(ctx, in)
}

func (c *ls2lc) LeaseRevoke(ctx context.Context, in *pb.LeaseRevokeRequest, opts ...grpc.CallOption) (*pb.LeaseRevokeResponse, error) {
	return c.leaseServer.LeaseRevoke(ctx, in)
}

func (c *ls2lc) LeaseKeepAlive(ctx context.Context, opts ...grpc.CallOption) (pb.Lease_LeaseKeepAliveClient, error) {
	cs := newPipeStream(ctx, func(ss chanServerStream) error {
		return c.leaseServer.LeaseKeepAlive(&ls2lcServerStream{ss})
	})
	return &ls2lcClientStream{cs}, nil
}

func (c *ls2lc) LeaseTimeToLive(ctx context.Context, in *pb.LeaseTimeToLiveRequest, opts ...grpc.CallOption) (*pb.LeaseTimeToLiveResponse, error) {
	return c.leaseServer.LeaseTimeToLive(ctx, in)
}

func (c *ls2lc) LeaseLeases(ctx context.Context, in *pb.LeaseLeasesRequest, opts ...grpc.CallOption) (*pb.LeaseLeasesResponse, error) {
	return c.leaseServer.LeaseLeases(ctx, in)
}

// ls2lcClientStream implements Lease_LeaseKeepAliveClient
type ls2lcClientStream struct{ chanClientStream }

// ls2lcServerStream implements Lease_LeaseKeepAliveServer
type ls2lcServerStream struct{ chanServerStream }

func (s *ls2lcClientStream) Send(rr *pb.LeaseKeepAliveRequest) error {
	return s.SendMsg(rr)
}
func (s *ls2lcClientStream) Recv() (*pb.LeaseKeepAliveResponse, error) {
	var v interface{}
	if err := s.RecvMsg(&v); err != nil {
		return nil, err
	}
	return v.(*pb.LeaseKeepAliveResponse), nil
}

func (s *ls2lcServerStream) Send(rr *pb.LeaseKeepAliveResponse) error {
	return s.SendMsg(rr)
}
func (s *ls2lcServerStream) Recv() (*pb.LeaseKeepAliveRequest, error) {
	var v interface{}
	if err := s.RecvMsg(&v); err != nil {
		return nil, err
	}
	return v.(*pb.LeaseKeepAliveRequest), nil
}
