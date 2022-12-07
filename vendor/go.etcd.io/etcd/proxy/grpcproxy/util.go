package grpcproxy

import (
	"context"

	"go.etcd.io/etcd/etcdserver/api/v3rpc/rpctypes"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func getAuthTokenFromClient(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ts, ok := md[rpctypes.TokenFieldNameGRPC]
		if ok {
			return ts[0]
		}
	}
	return ""
}

func withClientAuthToken(ctx, ctxWithToken context.Context) context.Context {
	token := getAuthTokenFromClient(ctxWithToken)
	if token != "" {
		ctx = context.WithValue(ctx, rpctypes.TokenFieldNameGRPC, token)
	}
	return ctx
}

type proxyTokenCredential struct {
	token string
}

func (cred *proxyTokenCredential) RequireTransportSecurity() bool {
	return false
}

func (cred *proxyTokenCredential) GetRequestMetadata(ctx context.Context, s ...string) (map[string]string, error) {
	return map[string]string{
		rpctypes.TokenFieldNameGRPC: cred.token,
	}, nil
}

func AuthUnaryClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	token := getAuthTokenFromClient(ctx)
	if token != "" {
		tokenCred := &proxyTokenCredential{token}
		opts = append(opts, grpc.PerRPCCredentials(tokenCred))
	}
	return invoker(ctx, method, req, reply, cc, opts...)
}

func AuthStreamClientInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	tokenif := ctx.Value(rpctypes.TokenFieldNameGRPC)
	if tokenif != nil {
		tokenCred := &proxyTokenCredential{tokenif.(string)}
		opts = append(opts, grpc.PerRPCCredentials(tokenCred))
	}
	return streamer(ctx, desc, cc, method, opts...)
}
