package podresources

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	podresourcesapi "k8s.io/kubernetes/pkg/kubelet/apis/podresources/v1alpha1"
	"k8s.io/kubernetes/pkg/kubelet/util"
)

// GetClient returns a client for the PodResourcesLister grpc service
func GetClient(socket string, connectionTimeout time.Duration, maxMsgSize int) (podresourcesapi.PodResourcesListerClient, *grpc.ClientConn, error) {
	addr, dialer, err := util.GetAddressAndDialer(socket)
	if err != nil {
		return nil, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), connectionTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithInsecure(), grpc.WithDialer(dialer), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(maxMsgSize)))
	if err != nil {
		return nil, nil, fmt.Errorf("error dialing socket %s: %v", socket, err)
	}
	return podresourcesapi.NewPodResourcesListerClient(conn), conn, nil
}
