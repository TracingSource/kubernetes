package audit

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apiserver/pkg/apis/audit/install"
	auditv1 "k8s.io/apiserver/pkg/apis/audit/v1"
)

func cleanup(t *testing.T, path string) {
	err := os.RemoveAll(path)
	if err != nil {
		t.Fatalf("Failed to clean up %v: %v", path, err)
	}
}

func TestCreateDefaultAuditLogPolicy(t *testing.T) {
	// make a tempdir
	tempDir, err := ioutil.TempDir("/tmp", "audit-test")
	if err != nil {
		t.Fatalf("could not create a tempdir: %v", err)
	}
	defer cleanup(t, tempDir)
	auditPolicyFile := filepath.Join(tempDir, "test.yaml")
	if err = CreateDefaultAuditLogPolicy(auditPolicyFile); err != nil {
		t.Fatalf("failed to create audit log policy: %v", err)
	}
	// turn the audit log back into a policy
	policyBytes, err := ioutil.ReadFile(auditPolicyFile)
	if err != nil {
		t.Fatalf("failed to read %v: %v", auditPolicyFile, err)
	}
	scheme := runtime.NewScheme()
	install.Install(scheme)
	codecs := serializer.NewCodecFactory(scheme)
	policy := auditv1.Policy{}
	err = runtime.DecodeInto(codecs.UniversalDecoder(), policyBytes, &policy)
	if err != nil {
		t.Fatalf("failed to decode written policy: %v", err)
	}
	if policy.Kind != "Policy" {
		t.Fatalf("did not decode policy properly")
	}
}
