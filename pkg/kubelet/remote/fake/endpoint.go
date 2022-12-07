// +build !windows

package fake

const (
	defaultUnixEndpoint = "unix:///tmp/kubelet_remote.sock"
)

// GenerateEndpoint generates a new unix socket server of grpc server.
func GenerateEndpoint() (string, error) {
	return defaultUnixEndpoint, nil
}
