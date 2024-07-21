# make generated_files 生成 openapi.GetOpenAPIDefinitions

初始 clone kubernetes 仓库并启动apiserver前, 需要先执行`make generated_files`, 生成如下文件.

- staging/src/k8s.io/code-generator/_examples/apiserver/openapi/zz_generated.openapi.go
- staging/src/k8s.io/sample-apiserver/pkg/generated/openapi/zz_generated.openapi.go
- staging/src/k8s.io/apiextensions-apiserver/pkg/generated/openapi/zz_generated.openapi.go

否则启动会报错

```log
app/server.go:553:3: undefined: "k8s.io/kubernetes/pkg/generated/openapi".GetOpenAPIDefinitions (exit status 2)
```

但是 1.17.2 版本将 stage 中的一些仓库(apimachinery, apiextensions-apiserver 等)独立出来后, 执行该命令会卡住.

```log
$ make generated_files DBG_MAKEFILE=1
Makefile:20: ***** starting Makefile for goal(s) "generated_files"
Makefile:21: ***** Mon Dec  5 16:30:01 CST 2022
make -f Makefile.generated_files generated_files CALLED_FROM_MAIN_MAKEFILE=1
make[1]: Entering directory `/home/k8s.io/kubernetes'
Makefile.generated_files:27: ***** starting Makefile.generated_files for goal(s) "generated_files"
Makefile.generated_files:28: ***** Mon Dec  5 16:30:01 CST 2022
Makefile.generated_files:93: ***** finding all *.go dirs
Makefile.generated_files:102: ***** finding all +k8s: tags
Makefile.generated_files:135: ***** finding all +k8s:deepcopy-gen tags
Makefile.generated_files:230: ***** finding all +k8s:defaulter-gen tags
Makefile.generated_files:332: ***** finding all +k8s:conversion-gen tags
Makefile.generated_files:451: ***** finding all +k8s:openapi-gen tags for KUBE

```

卡住不动了, ps一下看看卡在哪了.

```log
$ ps -ef | grep make
root      70152   1983  0 10:48 pts/0    00:00:00 make generated_files DBG_MAKEFILE=1
root      70154  70152  0 10:48 pts/0    00:00:00 make -f Makefile.generated_files generated_files CALLED_FROM_MAIN_MAKEFILE=1
root      70538   2418  0 10:49 pts/1    00:00:00 grep --color=auto make
$ ps -ef | grep 70154
root      70154  70152  0 10:48 pts/0    00:00:00 make -f Makefile.generated_files generated_files CALLED_FROM_MAIN_MAKEFILE=1
root      70522  70154  0 10:48 pts/0    00:00:00 /bin/bash -c grep --color=never -l '+k8s:openapi-gen='  | xargs -n1 dirname | LC_ALL=C sort -u
root      70540   2418  0 10:49 pts/1    00:00:00 grep --color=auto 70154
$ ps -ef | grep 70522
root      70522  70154  0 10:48 pts/0    00:00:00 /bin/bash -c grep --color=never -l '+k8s:openapi-gen='  | xargs -n1 dirname | LC_ALL=C sort -u
root      70523  70522  0 10:48 pts/0    00:00:00 grep --color=never -l +k8s:openapi-gen=
root      70524  70522  0 10:48 pts/0    00:00:00 xargs -n1 dirname
root      70525  70522  0 10:48 pts/0    00:00:00 sort -u
root      70542   2418  0 10:49 pts/1    00:00:00 grep --color=auto 70522
```

## 解决方法

后来在 1.16.2 版本的分支下, 输出了 makefile 中`ALL_GO_DIRS`变量, 该变量应该是包含所有拥有`.go`文件的目录列表, 使用空格分隔. 

对比发现 1.17.2 在拆分上述仓库后, 没有出现在该变量的内容中, 所以在`make generated_files`前, 要先把那几个独立的库拷到 stage 目录下, 然后再在 vendor 下添加好软链接.

```
mv ../apiextensions-apiserver   ./staging/src/k8s.io/apiextensions-apiserver
mv ../apimachinery              ./staging/src/k8s.io/apimachinery
mv ../apiserver                 ./staging/src/k8s.io/apiserver
mv ../client-go                 ./staging/src/k8s.io/client-go
mv ../kube-aggregator           ./staging/src/k8s.io/kube-aggregator
mv ../kubectl                   ./staging/src/k8s.io/kubectl
mv ../kube-api                  ./staging/src/k8s.io/api

cd vendor/k8s.io/

ln -s ../../staging/src/k8s.io/apiextensions-apiserver  ./apiextensions-apiserver
ln -s ../../staging/src/k8s.io/apimachinery             ./apimachinery           
ln -s ../../staging/src/k8s.io/apiserver                ./apiserver              
ln -s ../../staging/src/k8s.io/client-go                ./client-go              
ln -s ../../staging/src/k8s.io/kube-aggregator          ./kube-aggregator        
ln -s ../../staging/src/k8s.io/kubectl                  ./kubectl                
ln -s ../../staging/src/k8s.io/api                      ./api                    
```

需要注意的是, utils 包是直接放在 vendor 目录的, 没有用软链接.

```
mv ../kube-utils                ./vendor/k8s.io/utils
```

然后生成那3个文件后, 再移出来就可以了.

注意, 其中有一个文件是生成在 apiextensions-apiserver 工程中的, 独立出来之后注意提交.
