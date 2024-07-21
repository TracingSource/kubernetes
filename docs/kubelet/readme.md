
------

/var/run/dockershim.sock 是在启动 kubelet 的时候才会出现, docker 默认是没有这个文件的.

且在 systemctl stop kubelet 后并不会消失, 只在 kubeadm reset 才会删除.

这部分则需要从 kubeadm init 流程中去寻找.

------

kubelet 与 docker api 交互的终点在这个文件.

```
pkg/kubelet/remote/remote_runtime.go
```

通过建立与 /var/run/dockershim.sock 文件的连接进行通信.

主要操作容器, pause容器, 以及镜像, 这3类目标.

kubelet理论上会建立2个与dockershim的连接, 一个用来管理容器, 一个用来管理镜像.

------

## 维护Pod列表

kubelet 需要维护的 Pod 信息来源主要有2个, 一个是当前主机下的 manifests 表示的静态 Pod, 另一个就是来自 apiserver 中的常规 Pod.

kubelet 需要对比这2处来源"期望"的 Pod 状态, 与当前主机上正在运行的 Pod 状态, 让ta们保持一致, 实现"consile(调和)"的功能.

这两个来源在如下文件中被指定, 同时开始了处理逻辑.

```
pkg/kubelet/kubelet.go -> makePodSourceConfig()
```

### 静态Pod

对静态Pod列表(/etc/kubernetes/manifests目录)的监听操作在如下文件

```
pkg/kubelet/config/file.go -> NewSourceFile()
```

对应的事件处理入口在

```
pkg/util/config/config.go -> Mux.listen()
```

### apiserver

```
pkg/kubelet/config/apiserver.go -> NewSourceApiserver()
```

对应的事件处理入口在

```
pkg/util/config/config.go -> Mux.listen()
```

### 总结 

初始化 PodConfig.updates 通道.

pkg/kubelet/kubelet.go -> makePodSourceConfig()
|
├─  pkg/kubelet/config/file.go -> NewSourceFile()
├─  pkg/kubelet/config/apiserver.go -> NewSourceApiserver()
|
pkg/util/config/config.go -> Mux.listen()
pkg/kubelet/config/config.go -> podStorage.Merge() 将在 apiserver/manifests 发生的变动进行汇总, 写入 PodConfig.updates 通道.

kubelet 启动的主线流程, 层层传递 PodConfig.updates 通道.

pkg/kubelet/config/config.go -> PodConfig.Updates()
cmd/kubelet/app/server.go -> startKubelet()
pkg/kubelet/kubelet.go -> Kubelet.Run()
pkg/kubelet/kubelet.go -> Kubelet.syncLoop()
pkg/kubelet/kubelet.go -> Kubelet.syncLoopIteration(): 从 PodConfig.updates 通道中取出 ADD, UPDATE, DELETE, SET 等事件, 分别处理.
pkg/kubelet/pod/pod_manager.go -> basicManager.AddPod(), UpdatePod(), DeletePod()
