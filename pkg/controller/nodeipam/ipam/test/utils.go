package test

import (
	"net"
)

// MustParseCIDR returns the CIDR range parsed from s or panics if the string
// cannot be parsed.
func MustParseCIDR(s string) *net.IPNet {
	_, ret, err := net.ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return ret
}
