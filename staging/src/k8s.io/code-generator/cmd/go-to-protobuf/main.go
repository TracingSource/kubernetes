// go-to-protobuf generates a Protobuf IDL from a Go struct, respecting any
// existing IDL tags on the Go struct.
package main

import (
	goflag "flag"

	flag "github.com/spf13/pflag"
	"k8s.io/code-generator/cmd/go-to-protobuf/protobuf"
)

var g = protobuf.New()

func init() {
	g.BindFlags(flag.CommandLine)
	goflag.Set("logtostderr", "true")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
}

func main() {
	flag.Parse()
	protobuf.Run(g)
}
