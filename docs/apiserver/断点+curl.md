- [官方文档 apiserver API](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.19/#pod-v1-core)
    - timeoutSeconds

```
header_auth='Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IjVrMk5oYjRHUnZ2VVNoN0VXcTg0TjNvSzJFX3o4NjJFeUJwSjFoNWlZNG8ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJkZWZhdWx0LXRva2VuLXh6cGhyIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImRlZmF1bHQiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI0NzQ5NjcxZS0zNjgzLTQ4OGMtYmY3OS1jMTVkMjdjZjFiZGQiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06ZGVmYXVsdCJ9.VFLxhpOxFWyqC9A2vABTf8Z041aTB6GZtb5o7e9xMPcjmQ-FFXap_cLU_cKn8IQeuAYmV5TZs1RiZNMVfb-w42TSTRgvirLTaPkeeejFEWuPMyyr3YpHpiiE6aZPm2T_5FI-DCoyhQ-VHwwlIRXVMijsUjGQtRh2JZbk-9hP4ZO8Hva9INXM0TRD6icRZncvwQJmQfu_QNQiFJmhnhXmKRW8YM74RtVQMkukPo7WSEVDcyoeH9V8VRP3KdgsRvgfpChHEC2lcCnDLrfj2M5UGiy9DcIj_6ei8JCgsE94Hx_XNDSpuCzCVLSlAyOCvL7LqLG81Fj_0mdSvzk7WXDGbA'
header_json='Content-Type: application/json'

curl -k --max-time 3600 -H "$header_json" -H "$header_auth" 'https://127.0.0.1:16443/api/v1/namespaces/kube-system/pods'
curl -k --max-time 3600 -H "$header_json" -H "$header_auth" 'https://127.0.0.1:16443/api/v1/namespaces/kube-system/pods/kube-apiserver-k8s-master-01?timeoutSeconds=0'
curl -k --max-time 3600 -H "$header_json" -H "$header_auth" 'https://127.0.0.1:16443/api/v1/namespaces/kube-system/pods/kube-apiserver-k8s-master-01?timeout=3600s'

https://k8s-server-lb:8443/api/v1/nodes?fieldSelector=metadata.name%3Dk8s-master-01&limit=500&resourceVersion=0
```

查看 apiserver 所在容器/主机的 /var/log/ 目录下的文件列表.

```
curl -k --max-time 3600 -H "$header_json" -H "$header_auth" 'https://127.0.0.1:16443/logs/'
```

查看 /var/log/messages 文件的内容.

```
curl -k --max-time 3600 -H "$header_json" -H "$header_auth" 'https://127.0.0.1:16443/logs/messages'
```

查看指定 group/资源 所支持的操作类型(get, create, update)

```
curl -k --max-time 3600 -H "$header_json" -H "$header_auth" 'https://127.0.0.1:16443/api/v1'
curl -k --max-time 3600 -H "$header_json" -H "$header_auth" 'https://127.0.0.1:16443/apis/apiextensions.k8s.io/v1beta1'
```

查看目标类型的所有资源

```
curl -k --max-time 3600 -H "$header_json" -H "$header_auth" 'https://127.0.0.1:16443/api/v1/pods'
```
