package proxy

const (
	// KubeProxyConfigMap19 is the proxy ConfigMap manifest for Kubernetes 1.9 and above
	KubeProxyConfigMap19 = `
kind: ConfigMap
apiVersion: v1
metadata:
  name: {{ .ProxyConfigMap }}
  namespace: kube-system
  labels:
    app: kube-proxy
data:
  kubeconfig.conf: |-
    apiVersion: v1
    kind: Config
    clusters:
    - cluster:
        certificate-authority: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
        server: {{ .ControlPlaneEndpoint }}
      name: default
    contexts:
    - context:
        cluster: default
        namespace: default
        user: default
      name: default
    current-context: default
    users:
    - name: default
      user:
        tokenFile: /var/run/secrets/kubernetes.io/serviceaccount/token
  {{ .ProxyConfigMapKey }}: |-
{{ .ProxyConfig}}
`

	// KubeProxyDaemonSet19 is the proxy DaemonSet manifest for Kubernetes 1.9 and above
	KubeProxyDaemonSet19 = `
apiVersion: apps/v1
kind: DaemonSet
metadata:
  labels:
    k8s-app: kube-proxy
  name: kube-proxy
  namespace: kube-system
spec:
  selector:
    matchLabels:
      k8s-app: kube-proxy
  updateStrategy:
    type: RollingUpdate
  template:
    metadata:
      labels:
        k8s-app: kube-proxy
    spec:
      priorityClassName: system-node-critical
      containers:
      - name: kube-proxy
        image: {{ .Image }}
        imagePullPolicy: IfNotPresent
        command:
        - /usr/local/bin/kube-proxy
        - --config=/var/lib/kube-proxy/{{ .ProxyConfigMapKey }}
        - --hostname-override=$(NODE_NAME)
        securityContext:
          privileged: true
        volumeMounts:
        - mountPath: /var/lib/kube-proxy
          name: kube-proxy
        - mountPath: /run/xtables.lock
          name: xtables-lock
          readOnly: false
        - mountPath: /lib/modules
          name: lib-modules
          readOnly: true
        env:
          - name: NODE_NAME
            valueFrom:
              fieldRef:
                fieldPath: spec.nodeName
      hostNetwork: true
      serviceAccountName: kube-proxy
      volumes:
      - name: kube-proxy
        configMap:
          name: {{ .ProxyConfigMap }}
      - name: xtables-lock
        hostPath:
          path: /run/xtables.lock
          type: FileOrCreate
      - name: lib-modules
        hostPath:
          path: /lib/modules
      tolerations:
      - key: CriticalAddonsOnly
        operator: Exists
      - operator: Exists
      nodeSelector:
        beta.kubernetes.io/os: linux
`
)
