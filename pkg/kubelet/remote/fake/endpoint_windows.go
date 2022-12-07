// +build windows

package fake

import (
	"fmt"
	"net"
)

// GenerateEndpoint generates a new tcp endpoint of grpc server.
func GenerateEndpoint() (string, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return "", nil
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return "", err
	}

	defer l.Close()
	return fmt.Sprintf("tcp://127.0.0.1:%d", l.Addr().(*net.TCPAddr).Port), nil
}
