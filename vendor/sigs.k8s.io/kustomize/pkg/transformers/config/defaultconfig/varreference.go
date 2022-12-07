package defaultconfig

const (
	varReferenceFieldSpecs = `
varReference:
- path: spec/template/spec/initContainers/command
  kind: StatefulSet

- path: spec/template/spec/containers/command
  kind: StatefulSet

- path: spec/template/spec/initContainers/command
  kind: Deployment

- path: spec/template/spec/containers/command
  kind: Deployment

- path: spec/template/spec/initContainers/command
  kind: DaemonSet

- path: spec/template/spec/containers/command
  kind: DaemonSet

- path: spec/template/spec/containers/command
  kind: Job

- path: spec/jobTemplate/spec/template/spec/containers/command
  kind: CronJob

- path: spec/template/spec/initContainers/args
  kind: StatefulSet
 
- path: spec/template/spec/containers/args
  kind: StatefulSet

- path: spec/template/spec/initContainers/args
  kind: Deployment

- path: spec/template/spec/containers/args
  kind: Deployment

- path: spec/template/spec/initContainers/args
  kind: DaemonSet

- path: spec/template/spec/containers/args
  kind: DaemonSet

- path: spec/template/spec/containers/args
  kind: Job

- path: spec/jobTemplate/spec/template/spec/containers/args
  kind: CronJob

- path: spec/template/spec/initContainers/env/value
  kind: StatefulSet

- path: spec/template/spec/containers/env/value
  kind: StatefulSet

- path: spec/template/spec/initContainers/env/value
  kind: Deployment

- path: spec/template/spec/containers/env/value
  kind: Deployment

- path: spec/template/spec/initContainers/env/value
  kind: DaemonSet

- path: spec/template/spec/containers/env/value
  kind: DaemonSet

- path: spec/template/spec/containers/env/value
  kind: Job

- path: spec/jobTemplate/spec/template/spec/containers/env/value
  kind: CronJob

- path: spec/containers/command
  kind: Pod

- path: spec/containers/args
  kind: Pod

- path: spec/containers/env/value
  kind: Pod

- path: spec/initContainers/command
  kind: Pod

- path: spec/initContainers/args
  kind: Pod

- path: spec/initContainers/env/value
  kind: Pod

- path: spec/rules/host
  kind: Ingress

- path: spec/tls/hosts
  kind: Ingress

- path: spec/template/spec/containers/volumeMounts/mountPath
  kind: StatefulSet

- path: spec/template/spec/initContainers/volumeMounts/mountPath
  kind: StatefulSet

- path: spec/containers/volumeMounts/mountPath
  kind: Pod

- path: spec/initContainers/volumeMounts/mountPath
  kind: Pod

- path: spec/template/spec/containers/volumeMounts/mountPath
  kind: ReplicaSet

- path: spec/template/spec/initContainers/volumeMounts/mountPath
  kind: ReplicaSet

- path: spec/template/spec/containers/volumeMounts/mountPath
  kind: Job

- path: spec/template/spec/initContainers/volumeMounts/mountPath
  kind: Job

- path: spec/template/spec/containers/volumeMounts/mountPath
  kind: CronJob

- path: spec/template/spec/initContainers/volumeMounts/mountPath
  kind: CronJob

- path: spec/template/spec/containers/volumeMounts/mountPath
  kind: DaemonSet

- path: spec/template/spec/initContainers/volumeMounts/mountPath
  kind: DaemonSet

- path: spec/template/spec/containers/volumeMounts/mountPath
  kind: Deployment

- path: spec/template/spec/initContainers/volumeMounts/mountPath
  kind: Deployment

- path: metadata/labels
`
)
