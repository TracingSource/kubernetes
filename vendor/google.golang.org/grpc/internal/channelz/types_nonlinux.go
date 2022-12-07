// +build !linux appengine

/*
 *
 * Copyright 2018 gRPC authors.
 *

 *
 */

package channelz

import (
	"sync"

	"google.golang.org/grpc/grpclog"
)

var once sync.Once

// SocketOptionData defines the struct to hold socket option data, and related
// getter function to obtain info from fd.
// Windows OS doesn't support Socket Option
type SocketOptionData struct {
}

// Getsockopt defines the function to get socket options requested by channelz.
// It is to be passed to syscall.RawConn.Control().
// Windows OS doesn't support Socket Option
func (s *SocketOptionData) Getsockopt(fd uintptr) {
	once.Do(func() {
		grpclog.Warningln("Channelz: socket options are not supported on non-linux os and appengine.")
	})
}
