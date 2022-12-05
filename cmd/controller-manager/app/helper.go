/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package app

import (
	"fmt"
	"net/http"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog"
)

// WaitForAPIServer 测试与 apiserver 的连接(通过访问 https://localhost:6443/healthz).
//
// 	@param client: kube client
// 	@param timeout: 超时时间
//
// caller: 
// 	1. cmd/kube-controller-manager/app/controllermanager.go -> CreateControllerContext()
//
// WaitForAPIServer waits for the API Server's /healthz endpoint to report "ok" with timeout.
func WaitForAPIServer(client clientset.Interface, timeout time.Duration) error {
	var lastErr error

	err := wait.PollImmediate(time.Second, timeout, func() (bool, error) {
		healthStatus := 0
		result := client.Discovery().RESTClient().Get().AbsPath("/healthz").Do().StatusCode(&healthStatus)
		if result.Error() != nil {
			lastErr = fmt.Errorf("failed to get apiserver /healthz status: %v", result.Error())
			return false, nil
		}
		if healthStatus != http.StatusOK {
			content, _ := result.Raw()
			lastErr = fmt.Errorf("APIServer isn't healthy: %v", string(content))
			klog.Warningf("APIServer isn't healthy yet: %v. Waiting a little while.", string(content))
			return false, nil
		}

		return true, nil
	})

	if err != nil {
		return fmt.Errorf("%v: %v", err, lastErr)
	}

	return nil
}

// IsControllerEnabled 判断目标 controller 是否被启用
// (只有返回 true 的 controller 才会被启用).
//
// 	@param name: 目标 controller 名称, 如 namespace, deployment 等
// 	@param disabledByDefaultControllers: 这是一个固定的列表, 见
// 	cmd/kube-controller-manager/app/controllermanager.go -> ControllersDisabledByDefault()
//
// caller: 
// 	1. cmd/kube-controller-manager/app/controllermanager.go -> 
// 	ControllerContext.IsControllerEnabled() 
// 	kcm 在启动时, 启动各 controller. 每启动一种 controller, 就需要先判断一下其是否被启用.
//
// IsControllerEnabled check if a specified controller enabled or not.
func IsControllerEnabled(
	name string, disabledByDefaultControllers sets.String, controllers []string,
) bool {
	hasStar := false
	for _, ctrl := range controllers {
		if ctrl == name {
			return true
		}
		if ctrl == "-"+name {
			return false
		}
		if ctrl == "*" {
			hasStar = true
		}
	}
	// if we get here, there was no explicit choice
	if !hasStar {
		// nothing on by default
		return false
	}

	return !disabledByDefaultControllers.Has(name)
}
