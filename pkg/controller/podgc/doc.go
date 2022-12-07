// Package podgc contains a very simple pod "garbage collector" implementation,
// PodGCController, that runs in the controller manager. If the number of pods
// in terminated phases (right now either Failed or Succeeded) surpasses a
// configurable threshold, the controller will delete pods in terminated state
// until the system reaches the allowed threshold again. The PodGCController
// prioritizes pods to delete by sorting by creation timestamp and deleting the
// oldest objects first. The PodGCController will not delete non-terminated
// pods.
package podgc // import "k8s.io/kubernetes/pkg/controller/podgc"
