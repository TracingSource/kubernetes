package network

import (
	"k8s.io/kubernetes/test/e2e/framework"
	"k8s.io/kubernetes/test/e2e/network/scale"

	"github.com/onsi/ginkgo"
)

var _ = SIGDescribe("Loadbalancing: L7 Scalability", func() {
	defer ginkgo.GinkgoRecover()
	var (
		ns string
	)
	f := framework.NewDefaultFramework("ingress-scale")

	ginkgo.BeforeEach(func() {
		ns = f.Namespace.Name
	})

	ginkgo.Describe("GCE [Slow] [Serial] [Feature:IngressScale]", func() {
		var (
			scaleFramework *scale.IngressScaleFramework
		)

		ginkgo.BeforeEach(func() {
			framework.SkipUnlessProviderIs("gce", "gke")

			scaleFramework = scale.NewIngressScaleFramework(f.ClientSet, ns, framework.TestContext.CloudConfig)
			if err := scaleFramework.PrepareScaleTest(); err != nil {
				framework.Failf("Unexpected error while preparing ingress scale test: %v", err)
			}
		})

		ginkgo.AfterEach(func() {
			if errs := scaleFramework.CleanupScaleTest(); len(errs) != 0 {
				framework.Failf("Unexpected error while cleaning up ingress scale test: %v", errs)
			}
		})

		ginkgo.It("Creating and updating ingresses should happen promptly with small/medium/large amount of ingresses", func() {
			if errs := scaleFramework.RunScaleTest(); len(errs) != 0 {
				framework.Failf("Unexpected error while running ingress scale test: %v", errs)
			}

		})
	})
})
