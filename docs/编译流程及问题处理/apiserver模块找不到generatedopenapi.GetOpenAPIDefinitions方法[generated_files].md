# apiserver模块找不到generatedopenapi.GetOpenAPIDefinitions方法[generated_files]

参考文章

1. [kubernetes 代码编译包 undefined: "k8s.io/kubernetes/pkg/generated/openapi".GetOpenAPIDefinitions 的解决办法](https://blog.csdn.net/qq_21816375/article/details/84929541)
    - `make generated_files`

在vscode阅读kube-apiserver代码时, 发现[cmd/kube-apiserver/app/server.go](https://github.com/kubernetes/kubernetes/blob/v1.16.0/cmd/kube-apiserver/app/server.go#L425)部分的代码中出现了红色波浪线, 说找不到`GetOpenAPIDefinitions()`函数.

```go
import(
    ...
	generatedopenapi "k8s.io/kubernetes/pkg/generated/openapi"
    ...
)
...
	genericConfig.OpenAPIConfig = genericapiserver.DefaultOpenAPIConfig(generatedopenapi.GetOpenAPIDefinitions, openapinamer.NewDefinitionNamer(legacyscheme.Scheme, extensionsapiserver.Scheme, aggregatorscheme.Scheme))
```

进到`generatedopenapi`所引入的包中, 发现的确是没有这个函数...

![](https://gitee.com/generals-space/gitimg/raw/master/eab3343d1e061fe3beff2d199c8f9d20.png)

网上搜了下, 找到了参考文章1, 需要执行下`make generated_files`, 这会生成一个`zz_generated.openapi.go`文件.

![](https://gitee.com/generals-space/gitimg/raw/master/e032759729f3bc3d76df2ce4ec0f32e5.png)

> `kubernetes`工程不必放在GOPATH下执行`make generated_files`也可以.

注意: golang的版本应是 1.12, 否则会报如下错误.

```
make[1]: *** [_output/bin/deepcopy-gen] 错误 1
make: *** [generated_files] 错误 2
```
