package main

import (
	"k8s.io/kubernetes/cmd/kubeadm/app"
	"k8s.io/kubernetes/cmd/kubeadm/app/util"
)

func main() {
	util.CheckErr(app.Run())
}
