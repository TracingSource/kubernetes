// The kubelet binary is responsible for maintaining a set of containers on a particular host VM.
// It syncs data from both configuration file(s) as well as from a quorum of etcd servers.
// It then queries Docker to see what is currently running.  It synchronizes the configuration data,
// with the running set of containers by starting or stopping Docker containers.
package main

import (
	"math/rand"
	"os"
	"time"

	"k8s.io/component-base/logs"
	_ "k8s.io/component-base/metrics/prometheus/restclient"
	_ "k8s.io/component-base/metrics/prometheus/version" // for version metric registration
	"k8s.io/kubernetes/cmd/kubelet/app"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	command := app.NewKubeletCommand()
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
