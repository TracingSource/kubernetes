package node

import (
	"k8s.io/kubernetes/test/e2e/framework"
	e2esecurity "k8s.io/kubernetes/test/e2e/framework/security"

	"github.com/onsi/ginkgo"
)

var _ = SIGDescribe("AppArmor", func() {
	f := framework.NewDefaultFramework("apparmor")

	ginkgo.Context("load AppArmor profiles", func() {
		ginkgo.BeforeEach(func() {
			framework.SkipIfAppArmorNotSupported()
			e2esecurity.LoadAppArmorProfiles(f.Namespace.Name, f.ClientSet)
		})
		ginkgo.AfterEach(func() {
			if !ginkgo.CurrentGinkgoTestDescription().Failed {
				return
			}
			framework.LogFailedContainers(f.ClientSet, f.Namespace.Name, framework.Logf)
		})

		ginkgo.It("should enforce an AppArmor profile", func() {
			e2esecurity.CreateAppArmorTestPod(f.Namespace.Name, f.ClientSet, f.PodClient(), false, true)
		})

		ginkgo.It("can disable an AppArmor profile, using unconfined", func() {
			e2esecurity.CreateAppArmorTestPod(f.Namespace.Name, f.ClientSet, f.PodClient(), true, true)
		})
	})
})
