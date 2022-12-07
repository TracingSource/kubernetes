package latest

import (
	//Init the abac api package
	_ "k8s.io/kubernetes/pkg/apis/abac"
	_ "k8s.io/kubernetes/pkg/apis/abac/v0"
	_ "k8s.io/kubernetes/pkg/apis/abac/v1beta1"
)

// TODO: this file is totally wrong, it should look like other latest files.
// lavalamp is in the middle of fixing this code, so wait for the new way of doing things..
