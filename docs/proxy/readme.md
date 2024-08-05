参考文章

1. [k8s集群 iptables vs ipvs性能测试](https://www.flftuu.com/2019/11/01/k8s%E9%9B%86%E7%BE%A4-iptables-vs-ipvs%E6%80%A7%E8%83%BD%E6%B5%8B%E8%AF%95/)
2. [kubernetes网络相关总结](http://codemacro.com/2018/04/01/kube-network/)
    - 备份地址: [kubernetes网络相关总结](https://blog.csdn.net/dest_dest/article/details/80695734)
    - kuber 中 proxy 组件使用的 iptables 架构图, 值得一看.
3. [CHANGELOG-1.8](https://github.com/kubernetes/kubernetes/blob/master/CHANGELOG-1.8.md#kube-proxy-ipvs-mode) 
    - proxy的ipvs模型在1.8版本时由华为提交pr. 
4. [Implement IPVS-based in-cluster service load balancing](https://github.com/kubernetes/kubernetes/pull/46580)
    - 华为的 pr 地址
    - 提到此pr解决了什么问题, 与iptables模型相比有什么优势.
    - 但貌似直到1.17都没有实现GA.

##

`proxy`是kube组件中最外围的部分, 可以说ta就是一个简单的CRD. ta是用来监听集群中`service/endpoint`的变化, 然后将各service的ServiceIP与对应的pod资源的PodIP对应起来. 在ipvs模型中, 使用的ipvs的Virtual Server和Real Server. 其中ServiceIP作为VS, PodIP作为RS.

`proxy`只提供了Service到Pod映射, 并不提供跨宿主机通信的功能(这是由CNI插件来实现的部分). 所以在setup一个kuber集群时, 其余各组件(除了dns)使用的都是`hostNetwork`, 否则运行在Pod内的各组件是没有办法与其他宿主机通信的, 这也导致了在未部署CNI服务的时候, dns组件一直处于`ContainerCreating`的状态.

## iptables 存在的问题:

iptables 模式下对 node 节点的资源消耗明显大于 ipvs 模式

1. 规则顺序匹配延迟大
2. 访问 service 时需要遍历每条链知道匹配, 时间复杂度 O(N), 当规则数增加时, 匹配时间也增加
3. 规则更新延迟大
4. iptables 规则更新不是增量式的, 每更新一条规则, 都会把全部规则加载刷新一遍
5. 规则数大时, 会出现 kernel lock
6. svc 数增加到 5000 时, 会频繁出现"Another app is currently holding the xtables lock. Stopped waiting after 5s", 导致规则更新延迟变大, `kube-proxy`定期同步时也会因为超时导致 CrashLoopBackOff
