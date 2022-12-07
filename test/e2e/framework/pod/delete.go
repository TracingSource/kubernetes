package pod

import (
	"fmt"
	"time"

	"github.com/onsi/ginkgo"

	v1 "k8s.io/api/core/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	clientset "k8s.io/client-go/kubernetes"
	e2elog "k8s.io/kubernetes/test/e2e/framework/log"
)

const (
	// PodDeleteTimeout is how long to wait for a pod to be deleted.
	PodDeleteTimeout = 5 * time.Minute
)

// DeletePodOrFail deletes the pod of the specified namespace and name.
func DeletePodOrFail(c clientset.Interface, ns, name string) {
	ginkgo.By(fmt.Sprintf("Deleting pod %s in namespace %s", name, ns))
	err := c.CoreV1().Pods(ns).Delete(name, nil)
	expectNoError(err, "failed to delete pod %s in namespace %s", name, ns)
}

// DeletePodWithWait deletes the passed-in pod and waits for the pod to be terminated. Resilient to the pod
// not existing.
func DeletePodWithWait(c clientset.Interface, pod *v1.Pod) error {
	if pod == nil {
		return nil
	}
	return DeletePodWithWaitByName(c, pod.GetName(), pod.GetNamespace())
}

// DeletePodWithWaitByName deletes the named and namespaced pod and waits for the pod to be terminated. Resilient to the pod
// not existing.
func DeletePodWithWaitByName(c clientset.Interface, podName, podNamespace string) error {
	e2elog.Logf("Deleting pod %q in namespace %q", podName, podNamespace)
	err := c.CoreV1().Pods(podNamespace).Delete(podName, nil)
	if err != nil {
		if apierrs.IsNotFound(err) {
			return nil // assume pod was already deleted
		}
		return fmt.Errorf("pod Delete API error: %v", err)
	}
	e2elog.Logf("Wait up to %v for pod %q to be fully deleted", PodDeleteTimeout, podName)
	err = WaitForPodNotFoundInNamespace(c, podName, podNamespace, PodDeleteTimeout)
	if err != nil {
		return fmt.Errorf("pod %q was not deleted: %v", podName, err)
	}
	return nil
}
