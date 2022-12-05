# 认证+鉴权 serviceaccount secret token jwt base64

认证

```
staging/src/k8s.io/apiserver/pkg/endpoints/filters/authentication.go -> WithAuthentication()
```

鉴权

```
staging/src/k8s.io/apiserver/pkg/endpoints/filters/authorization.go -> WithAuthorization()
```

将从请求头中取出的token 数据 , 与本地存储的数据进行比对, 以转换成对应的 User/Group 对象.

```
staging/src/k8s.io/apiserver/pkg/authentication/token/cache/cached_token_authenticator.go
```

一个 serviceAccount 所绑定的 secret 资源中的 token 格式大概如下

```
header_auth='Authorization: Bearer eyJhbGciOiJSUzI1NiIsImtpZCI6IjVrMk5oYjRHUnZ2VVNoN0VXcTg0TjNvSzJFX3o4NjJFeUJwSjFoNWlZNG8ifQ.eyJpc3MiOiJrdWJlcm5ldGVzL3NlcnZpY2VhY2NvdW50Iiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9uYW1lc3BhY2UiOiJrdWJlLXN5c3RlbSIsImt1YmVybmV0ZXMuaW8vc2VydmljZWFjY291bnQvc2VjcmV0Lm5hbWUiOiJkZWZhdWx0LXRva2VuLXh6cGhyIiwia3ViZXJuZXRlcy5pby9zZXJ2aWNlYWNjb3VudC9zZXJ2aWNlLWFjY291bnQubmFtZSI6ImRlZmF1bHQiLCJrdWJlcm5ldGVzLmlvL3NlcnZpY2VhY2NvdW50L3NlcnZpY2UtYWNjb3VudC51aWQiOiI0NzQ5NjcxZS0zNjgzLTQ4OGMtYmY3OS1jMTVkMjdjZjFiZGQiLCJzdWIiOiJzeXN0ZW06c2VydmljZWFjY291bnQ6a3ViZS1zeXN0ZW06ZGVmYXVsdCJ9.VFLxhpOxFWyqC9A2vABTf8Z041aTB6GZtb5o7e9xMPcjmQ-FFXap_cLU_cKn8IQeuAYmV5TZs1RiZNMVfb-w42TSTRgvirLTaPkeeejFEWuPMyyr3YpHpiiE6aZPm2T_5FI-DCoyhQ-VHwwlIRXVMijsUjGQtRh2JZbk-9hP4ZO8Hva9INXM0TRD6icRZncvwQJmQfu_QNQiFJmhnhXmKRW8YM74RtVQMkukPo7WSEVDcyoeH9V8VRP3KdgsRvgfpChHEC2lcCnDLrfj2M5UGiy9DcIj_6ei8JCgsE94Hx_XNDSpuCzCVLSlAyOCvL7LqLG81Fj_0mdSvzk7WXDGbA'

```

`Bearer`后面的字符串可以用`.`进行分隔, 得到3个 part.

对`part[1]`作base64解码可以得到如下内容

```json
{
    "iss":"kubernetes/serviceaccount",
    "kubernetes.io/serviceaccount/namespace":"kube-system",
    "kubernetes.io/serviceaccount/secret.name":"default-token-xzphr",
    "kubernetes.io/serviceaccount/service-account.name":"default",
    "kubernetes.io/serviceaccount/service-account.uid":"4749671e-3683-488c-bf79-c15d27cf1bdd",
    "sub":"system:serviceaccount:kube-system:default"
}
```
