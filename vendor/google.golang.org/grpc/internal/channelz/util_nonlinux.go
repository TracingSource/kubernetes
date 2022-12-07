// +build !linux appengine

/*
 *
 * Copyright 2018 gRPC authors.
 *

 *
 */

package channelz

// GetSocketOption gets the socket option info of the conn.
func GetSocketOption(c interface{}) *SocketOptionData {
	return nil
}
