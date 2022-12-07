package monitoring

import (
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/framework/metrics"
	e2enode "k8s.io/kubernetes/test/e2e/framework/node"
	instrumentation "k8s.io/kubernetes/test/e2e/instrumentation/common"

	gin "github.com/onsi/ginkgo"
	gom "github.com/onsi/gomega"
)

var _ = instrumentation.SIGDescribe("MetricsGrabber", func() {
	f := framework.NewDefaultFramework("metrics-grabber")
	var c, ec clientset.Interface
	var grabber *metrics.Grabber
	gin.BeforeEach(func() {
		var err error
		c = f.ClientSet
		ec = f.KubemarkExternalClusterClientSet
		framework.ExpectNoError(err)
		grabber, err = metrics.NewMetricsGrabber(c, ec, true, true, true, true, true)
		framework.ExpectNoError(err)
	})

	gin.It("should grab all metrics from API server.", func() {
		gin.By("Connecting to /metrics endpoint")
		response, err := grabber.GrabFromAPIServer()
		framework.ExpectNoError(err)
		gom.Expect(response).NotTo(gom.BeEmpty())
	})

	gin.It("should grab all metrics from a Kubelet.", func() {
		gin.By("Proxying to Node through the API server")
		node, err := e2enode.GetRandomReadySchedulableNode(f.ClientSet)
		framework.ExpectNoError(err)
		response, err := grabber.GrabFromKubelet(node.Name)
		framework.ExpectNoError(err)
		gom.Expect(response).NotTo(gom.BeEmpty())
	})

	gin.It("should grab all metrics from a Scheduler.", func() {
		gin.By("Proxying to Pod through the API server")
		// Check if master Node is registered
		nodes, err := c.CoreV1().Nodes().List(metav1.ListOptions{})
		framework.ExpectNoError(err)

		var masterRegistered = false
		for _, node := range nodes.Items {
			if strings.HasSuffix(node.Name, "master") {
				masterRegistered = true
			}
		}
		if !masterRegistered {
			framework.Logf("Master is node api.Registry. Skipping testing Scheduler metrics.")
			return
		}
		response, err := grabber.GrabFromScheduler()
		framework.ExpectNoError(err)
		gom.Expect(response).NotTo(gom.BeEmpty())
	})

	gin.It("should grab all metrics from a ControllerManager.", func() {
		gin.By("Proxying to Pod through the API server")
		// Check if master Node is registered
		nodes, err := c.CoreV1().Nodes().List(metav1.ListOptions{})
		framework.ExpectNoError(err)

		var masterRegistered = false
		for _, node := range nodes.Items {
			if strings.HasSuffix(node.Name, "master") {
				masterRegistered = true
			}
		}
		if !masterRegistered {
			framework.Logf("Master is node api.Registry. Skipping testing ControllerManager metrics.")
			return
		}
		response, err := grabber.GrabFromControllerManager()
		framework.ExpectNoError(err)
		gom.Expect(response).NotTo(gom.BeEmpty())
	})
})
