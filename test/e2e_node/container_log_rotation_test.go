package e2e_node

import (
	"time"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/features"
	kubeletconfig "k8s.io/kubernetes/pkg/kubelet/apis/config"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	kubelogs "k8s.io/kubernetes/pkg/kubelet/logs"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
	"k8s.io/kubernetes/test/e2e/framework"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const (
	testContainerLogMaxFiles    = 3
	testContainerLogMaxSize     = "40Ki"
	rotationPollInterval        = 5 * time.Second
	rotationEventuallyTimeout   = 3 * time.Minute
	rotationConsistentlyTimeout = 2 * time.Minute
)

var _ = framework.KubeDescribe("ContainerLogRotation [Slow] [Serial] [Disruptive]", func() {
	f := framework.NewDefaultFramework("container-log-rotation-test")
	ginkgo.Context("when a container generates a lot of log", func() {
		ginkgo.BeforeEach(func() {
			if framework.TestContext.ContainerRuntime != kubetypes.RemoteContainerRuntime {
				framework.Skipf("Skipping ContainerLogRotation test since the container runtime is not remote")
			}
		})

		tempSetCurrentKubeletConfig(f, func(initialConfig *kubeletconfig.KubeletConfiguration) {
			initialConfig.FeatureGates[string(features.CRIContainerLogRotation)] = true
			initialConfig.ContainerLogMaxFiles = testContainerLogMaxFiles
			initialConfig.ContainerLogMaxSize = testContainerLogMaxSize
		})

		ginkgo.It("should be rotated and limited to a fixed amount of files", func() {
			ginkgo.By("create log container")
			pod := &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "test-container-log-rotation",
				},
				Spec: v1.PodSpec{
					RestartPolicy: v1.RestartPolicyNever,
					Containers: []v1.Container{
						{
							Name:  "log-container",
							Image: busyboxImage,
							Command: []string{
								"sh",
								"-c",
								// ~12Kb/s. Exceeding 40Kb in 4 seconds. Log rotation period is 10 seconds.
								"while true; do echo hello world; sleep 0.001; done;",
							},
						},
					},
				},
			}
			pod = f.PodClient().CreateSync(pod)
			ginkgo.By("get container log path")
			framework.ExpectEqual(len(pod.Status.ContainerStatuses), 1)
			id := kubecontainer.ParseContainerID(pod.Status.ContainerStatuses[0].ContainerID).ID
			r, _, err := getCRIClient()
			framework.ExpectNoError(err)
			status, err := r.ContainerStatus(id)
			framework.ExpectNoError(err)
			logPath := status.GetLogPath()
			ginkgo.By("wait for container log being rotated to max file limit")
			gomega.Eventually(func() (int, error) {
				logs, err := kubelogs.GetAllLogs(logPath)
				if err != nil {
					return 0, err
				}
				return len(logs), nil
			}, rotationEventuallyTimeout, rotationPollInterval).Should(gomega.Equal(testContainerLogMaxFiles), "should eventually rotate to max file limit")
			ginkgo.By("make sure container log number won't exceed max file limit")
			gomega.Consistently(func() (int, error) {
				logs, err := kubelogs.GetAllLogs(logPath)
				if err != nil {
					return 0, err
				}
				return len(logs), nil
			}, rotationConsistentlyTimeout, rotationPollInterval).Should(gomega.BeNumerically("<=", testContainerLogMaxFiles), "should never exceed max file limit")
		})
	})
})
