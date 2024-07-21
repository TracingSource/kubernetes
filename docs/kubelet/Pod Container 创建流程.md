
```log
pkg/kubelet/kubelet.go
Kubelet.HandlePodAdditions(), Kubelet.HandlePodUpdates(), Kubelet.HandlePodReconcile()
|   HandlePodAdditions() 通过调用 Kubelet.dispatchWork() 调用 Kubelet.podWorkers.UpdatePod(), 
|   所以 UpdatePod() 是最终用于部署 Pod 工作的方法.
|
├─  pkg/kubelet/pod_workers.go
    podWorkers.UpdatePod()
    |
    ├─  podWorkers.managePodLoop():
        |
        ├─  podWorkers.syncPodFn() (实际是 Kubelet.syncPod() 方法)
            |
            ├─  pkg/kubelet/kubelet.go
            |   Kubelet.containerManager.NewPodContainerManager()
            ├─  pkg/kubelet/cm/pod_container_manager_linux.go
            |   podContainerManagerImpl.EnsureExists()
            |   |   为目标 Pod 创建 cgroup 层级.
            |   |
            |   ├─  podContainerManagerImpl.cgroupManager.Create()
            |
            ├─  Kubelet.containerRuntime.SyncPod()
                pkg/kubelet/kuberuntime/kuberuntime_manager.go
                kubeGenericRuntimeManager.SyncPod()
                |
                ├─  kubeGenericRuntimeManager.computePodActions()
                |   1. 计算 sandbox 和 container 的差异变化.
                |
                ├─  kubeGenericRuntimeManager.killPodWithSyncResult()
                |   2. 如果 sandbox 发生变化, 对应的 container 需要删除重建.
                |
                ├─  
                |
                ├─  kubeGenericRuntimeManager.createPodSandbox()
                |   4. 如果有必要, 为 Pod 创建 sandbox 容器.
                |
                ├─  kubeGenericRuntimeManager.startContainer
                |   |   启动启动常规的 container, 用于初始化的 initContainers,
                |   |   和临时调试的 ephemeralContainers
                |   |
                |   ├─  pkg/kubelet/kuberuntime/kuberuntime_container.go
                |   |   kubeGenericRuntimeManager.startContainer()
                |   |   |
                |   |   ├─  kubeGenericRuntimeManager.imagePuller.EnsureImageExists()
                |   |   |   1. 拉取镜像
                |   |   ├─  pkg/kubelet/container/ref.go
                |   |   |   GenerateContainerRef()
                |   |   |
                |   |   ├─  pkg/kubelet/kuberuntime/kuberuntime_container.go
                |   |   |   kubeGenericRuntimeManager.runtimeService.CreateContainer()
                |   |   |   2. 创建容器 一直忘记 docker 还有一个 create 子命令, 可以只创建而不启动.
                |   |   |
                |   |   ├─  pkg/kubelet/kuberuntime/kuberuntime_container.go
                |   |   |   kubeGenericRuntimeManager.runtimeService.StartContainer()
                |   |   |   |   3. 启动容器
                |   |   |   |
                |   |   |   ├─  pkg/kubelet/remote/remote_runtime.go
                |   |   |       RemoteRuntimeService.StartContainer()
                |   |   |       |
                |   |   |       ├─  RemoteRuntimeService.runtimeClient.StartContainer()
                |   |   |           这里的 runtimeClient 是连接 /var/run/dockershim.sock 的
                |   |   |           grpc 客户端对象.
                |   |   |
                |   |   ├─  pkg/kubelet/kuberuntime/kuberuntime_container.go
                |   |   |   kubeGenericRuntimeManager.runner.Run()
                |   |   |   4. 执行 PostStart 钩子函数
                |   |
                |   |
                |
                ├─
                ├─
```
