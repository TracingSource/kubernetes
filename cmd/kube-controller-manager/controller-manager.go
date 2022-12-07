// controller manager 用于监控 replication controllers, 并创建相关的pod, 使其达到rc期望的状态.
// 提供了API监控新的controller, 与增删pod的操作.
//
// The controller manager is responsible for monitoring replication controllers,
// and creating corresponding pods to achieve the desired state. 
// It uses the API to listen for new controllers and to create/delete pods.
package main

import (
	"math/rand"
	"os"
	"time"

	"k8s.io/component-base/logs"
	_ "k8s.io/component-base/metrics/prometheus/clientgo" // load all the prometheus client-go plugin
	_ "k8s.io/component-base/metrics/prometheus/version"  // for version metric registration
	"k8s.io/kubernetes/cmd/kube-controller-manager/app"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	command := app.NewControllerManagerCommand()

	// TODO: once we switch everything over to Cobra commands, we can go back to calling
	// utilflag.InitFlags() (by removing its pflag.Parse() call). For now, we have to set the
	// normalize func and add the go flag set by hand.
	// utilflag.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()

	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
