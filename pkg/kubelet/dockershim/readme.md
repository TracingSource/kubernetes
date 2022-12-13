# dockershim

kubernetes 最开始出现时, 是与 docker 强绑定的(当时没有其他容器化实现), kubelet 与 dockerd 直接通信.

后来才出现了 docker 以外的其他 runtime, 如 runv, rkt. 

2016年, kubernetes 官方发布了 cri 接口规范, 规范所有运行时接口. 但此时 docker 也发布了 swarm, 进行容器编排. 一个由下往下, 一个由下向上, 都向对方发起正义的背刺😂.

docker 没有理会这个 cri, kubernetes 官方只能自己写了个`dockershim`包, 给 docker 服务提供了 cri 适配. 

------

dockershim 是一个 GRPC 服务, ta 监听 /var/run/dockershim.sock 接口(类似于 http 端口), kubelet 在启动时会同时启动.

**protobuf**

[cri-api](https://github.com/kubernetes/cri-api)工程定义了 dockershim 提供的函数原型(protobuf).

**Server**

kubernetes:pkg/kubelet/dockershim/docker_service.go -> dockerService{} 定义了这个 grpc 的服务端处理函数.

dockerService 中包含一个 client 成员对象, 这个对象是 dockershim 服务与 dockerd 服务(/var/run/docker.sock)通信的客户端, 在初始化时就会与 dockerd 建立连接.

**client**

kubernetes:pkg/kubelet/remote/remote_runtime.go -> [RemoteRuntimeService{}, RemoteImageService{}] 这2个结构体, 则定义了 grpc 的客户端函数. 

kubelet 只通过这2个结构体与 dockerd 服务通信.

```
        |      GRPC client     |                   GRPC server                  |

        ┌ RemoteRuntimeService ┐
kubelet ┤                      ├──> dockershim.sock ──> dockerService ──> client ──> docker.sock ──> dockerd
        └  RemoteImageService  ┘
```
