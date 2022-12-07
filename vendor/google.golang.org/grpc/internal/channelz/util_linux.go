// +build linux,!appengine

/*
 *
 * Copyright 2018 gRPC authors.
 *

 *
 */

package channelz

import (
	"syscall"
)

// GetSocketOption gets the socket option info of the conn.
func GetSocketOption(socket interface{}) *SocketOptionData {
	c, ok := socket.(syscall.Conn)
	if !ok {
		return nil
	}
	data := &SocketOptionData{}
	if rawConn, err := c.SyscallConn(); err == nil {
		rawConn.Control(data.Getsockopt)
		return data
	}
	return nil
}
