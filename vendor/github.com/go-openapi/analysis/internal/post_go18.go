// +build go1.8

package internal

import "net/url"

// PathUnescape provides url.PathUnescape(), with seamless
// go version support for pre-go1.8
//
// TODO: this function is currently defined in go-openapi/swag,
// but unexported. We might chose to export it, or simple phase
// out pre-go1.8 support.
func PathUnescape(path string) (string, error) {
	return url.PathUnescape(path)
}
