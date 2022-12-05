参考文章

1. [APIServer深入解读](https://www.zybuluo.com/dujun/note/64868)
    1. APIServer职能
    2. APIServer启动过程
    3. 资源对象的RESTful接口
    4. API操作的原子性
    5. APIServer总体架构小结
2. [Kubernetes 源码分析之 Apiserver](https://www.jianshu.com/p/1dad72e45230)
    - 这篇文章中关于 Cacher 部分的流程还没有搞清楚.

APIGroupInfo{
    Group:  apps, batch, crd, crd/status
    Version: v1, v1beta1
}
根据 APIGroupInfo.Version 进行细分, 每个 Version 都可生成一个 APIGroupVersion 对象.
APIGroupVersion {

}

遍历目标资源(pod, pod/log等)支持的所有操作(get/list等), 将不同的操作注册为不同的路由(Create对应Post, Update对应Put等), 并与后端 Storage 的不同方法绑定.

比如 GET 请求直接转发到 storage 的 Get() 方法, POST 请求则转发到 storage 的 Create() 方法.

```
staging/src/k8s.io/apiserver/pkg/endpoints/installer_registerResourceHandlers.go
```

实际的 handler 操作, 在下面目录中, 包含 get, put, post 等操作.

```
staging/src/k8s.io/apiserver/pkg/endpoints/handlers
```

## 权限认证, 限流等中间件

在如下函数中被注册.

```
staging/src/k8s.io/apiserver/pkg/server/config.go -> DefaultBuildHandlerChain()
```
