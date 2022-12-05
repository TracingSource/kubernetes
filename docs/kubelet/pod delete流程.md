pkg/kubelet/kubelet.go -> Kubelet.syncPod()
    判断Pod的事件类型, 如果为 kill 事件, 则调用 killPod() 进行删除.

pkg/kubelet/kubelet_pods.go -> Kubelet.killPod()

pkg/kubelet/kuberuntime/kuberuntime_manager.go -> kubeGenericRuntimeManager.KillPod() -> kubeGenericRuntimeManager. killPodWithSyncResult()
    kill掉目标Pod中的所有 containers, 并同步获取删除结果.
    完成后还要调用 gc 服务回收 sandbox 容器.
pkg/kubelet/kuberuntime/kuberuntime_container.go -> kubeGenericRuntimeManager.killContainersWithSyncResult()
    并行删除 containers, 并将结果通过 channel 收集并返回.

pkg/kubelet/kuberuntime/kuberuntime_container.go -> kubeGenericRuntimeManager.killContainer()
    kill 单个 docker 容器, 且在此函数中将执行 preStop 的钩子.

pkg/kubelet/remote/remote_runtime.go -> RemoteRuntimeService.StopContainer()
    其实这里已经是调用 docker-shim 接口了.

