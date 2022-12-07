// +build appengine

/*
 *
 * Copyright 2018 gRPC authors.
 *

 *
 */

package internal

import (
	"net"
)

// WrapSyscallConn returns newConn on appengine.
func WrapSyscallConn(rawConn, newConn net.Conn) net.Conn {
	return newConn
}
