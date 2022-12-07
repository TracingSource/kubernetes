package config

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	apiequality "k8s.io/apimachinery/pkg/api/equality"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
)

func TestExtractFromBadDataFile(t *testing.T) {
	dirName, err := mkTempDir("file-test")
	if err != nil {
		t.Fatalf("unable to create temp dir: %v", err)
	}
	defer removeAll(dirName, t)

	fileName := filepath.Join(dirName, "test_pod_config")
	err = ioutil.WriteFile(fileName, []byte{1, 2, 3}, 0555)
	if err != nil {
		t.Fatalf("unable to write test file %#v", err)
	}

	ch := make(chan interface{}, 1)
	lw := newSourceFile(fileName, "localhost", time.Millisecond, ch)
	err = lw.listConfig()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	expectEmptyChannel(t, ch)
}

func TestExtractFromEmptyDir(t *testing.T) {
	dirName, err := mkTempDir("file-test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer removeAll(dirName, t)

	ch := make(chan interface{}, 1)
	lw := newSourceFile(dirName, "localhost", time.Millisecond, ch)
	err = lw.listConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	update, ok := (<-ch).(kubetypes.PodUpdate)
	if !ok {
		t.Fatalf("unexpected type: %#v", update)
	}
	expected := CreatePodUpdate(kubetypes.SET, kubetypes.FileSource)
	if !apiequality.Semantic.DeepEqual(expected, update) {
		t.Fatalf("expected %#v, got %#v", expected, update)
	}
}

func mkTempDir(prefix string) (string, error) {
	return ioutil.TempDir(os.TempDir(), prefix)
}

func removeAll(dir string, t *testing.T) {
	if err := os.RemoveAll(dir); err != nil {
		t.Fatalf("unable to remove dir %s: %v", dir, err)
	}
}
