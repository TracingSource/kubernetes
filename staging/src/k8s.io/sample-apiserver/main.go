package main

import (
	"flag"
	"os"

	"k8s.io/klog"

	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/logs"
	"k8s.io/sample-apiserver/pkg/cmd/server"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	stopCh := genericapiserver.SetupSignalHandler()
	options := server.NewWardleServerOptions(os.Stdout, os.Stderr)
	cmd := server.NewCommandStartWardleServer(options, stopCh)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Execute(); err != nil {
		klog.Fatal(err)
	}
}
