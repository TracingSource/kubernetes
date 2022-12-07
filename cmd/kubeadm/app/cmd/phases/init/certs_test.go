package phases

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	"k8s.io/kubernetes/cmd/kubeadm/app/cmd/phases/workflow"
	"k8s.io/kubernetes/cmd/kubeadm/app/phases/certs"
	certstestutil "k8s.io/kubernetes/cmd/kubeadm/app/util/certs"
	"k8s.io/kubernetes/cmd/kubeadm/app/util/pkiutil"
	testutil "k8s.io/kubernetes/cmd/kubeadm/test"
)

type testCertsData struct {
	testInitData
	cfg *kubeadmapi.InitConfiguration
}

func (t *testCertsData) Cfg() *kubeadmapi.InitConfiguration { return t.cfg }
func (t *testCertsData) ExternalCA() bool                   { return false }
func (t *testCertsData) CertificateDir() string             { return t.cfg.CertificatesDir }
func (t *testCertsData) CertificateWriteDir() string        { return t.cfg.CertificatesDir }

func TestCertsWithCSRs(t *testing.T) {
	csrDir := testutil.SetupTempDir(t)
	defer os.RemoveAll(csrDir)
	certDir := testutil.SetupTempDir(t)
	defer os.RemoveAll(certDir)
	cert := &certs.KubeadmCertAPIServer

	certsData := &testCertsData{
		cfg: testutil.GetDefaultInternalConfig(t),
	}
	certsData.cfg.CertificatesDir = certDir

	// global vars
	csrOnly = true
	csrDir = certDir
	defer func() {
		csrOnly = false
	}()

	phase := NewCertsPhase()
	// find the api cert phase
	var apiServerPhase *workflow.Phase
	for _, phase := range phase.Phases {
		if phase.Name == cert.Name {
			apiServerPhase = &phase
			break
		}
	}

	if apiServerPhase == nil {
		t.Fatalf("couldn't find apiserver phase")
	}

	err := apiServerPhase.Run(certsData)
	if err != nil {
		t.Fatalf("couldn't run API server phase: %v", err)
	}

	if _, _, err := pkiutil.TryLoadCSRAndKeyFromDisk(csrDir, cert.BaseName); err != nil {
		t.Fatalf("couldn't load certificate %q: %v", cert.BaseName, err)
	}
}

func TestCreateSparseCerts(t *testing.T) {
	for _, test := range certstestutil.GetSparseCertTestCases(t) {
		t.Run(test.Name, func(t *testing.T) {
			tmpdir := testutil.SetupTempDir(t)
			defer os.RemoveAll(tmpdir)

			certstestutil.WritePKIFiles(t, tmpdir, test.Files)

			r := workflow.NewRunner()
			r.AppendPhase(NewCertsPhase())
			r.SetDataInitializer(func(*cobra.Command, []string) (workflow.RunData, error) {
				certsData := &testCertsData{
					cfg: testutil.GetDefaultInternalConfig(t),
				}
				certsData.cfg.CertificatesDir = tmpdir
				return certsData, nil
			})

			if err := r.Run([]string{}); (err != nil) != test.ExpectError {
				t.Fatalf("expected error to be %t, got %t (%v)", test.ExpectError, (err != nil), err)
			}
		})
	}
}
