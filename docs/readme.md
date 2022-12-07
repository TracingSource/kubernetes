v1.17.2

源码编译生成的可执行文件在`_output/bin`目录下.

版本1.16需要使用go1.12编译, 使用1.13会失败. kuber各版本与对应的golang版本表格可见[Building Kubernetes on a local OS/shell environment](https://github.com/kubernetes/community/blob/master/contributors/devel/development.md#go)

编译单独组件

```
make WHAT=cmd/kube-controller-manager
make WHAT=cmd/kubectl
make WHAT=cmd/kube-apiserver
```

git clone https://gitee.com/skeyes/kubernetes.git
git clone https://gitee.com/skeyes/apimachinery.git
git clone https://gitee.com/skeyes/kube-api.git
git clone https://gitee.com/skeyes/kube-utils.git
git clone https://gitee.com/skeyes/kubectl.git
git clone https://gitee.com/skeyes/apiserver.git
git clone https://gitee.com/skeyes/kube-aggregator.git
git clone https://gitee.com/skeyes/apiextensions-apiserver.git
git clone https://gitee.com/skeyes/client-go.git

