package checkpoint

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager"
)

// TestWriteLoadDeletePods validates all combinations of write, load, and delete
func TestWriteLoadDeletePods(t *testing.T) {
	testPods := []struct {
		pod     *v1.Pod
		written bool
	}{
		{
			pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "Foo",
					Annotations: map[string]string{core.BootstrapCheckpointAnnotationKey: "true"},
					UID:         "1",
				},
			},
			written: true,
		},
		{
			pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:        "Foo2",
					Annotations: map[string]string{core.BootstrapCheckpointAnnotationKey: "true"},
					UID:         "2",
				},
			},
			written: true,
		},
		{
			pod: &v1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name: "Bar",
					UID:  "3",
				},
			},
			written: false,
		},
	}

	dir, err := ioutil.TempDir("", "checkpoint")
	if err != nil {
		t.Errorf("Failed to allocate temp directory for TestWriteLoadDeletePods error=%v", err)
	}
	defer os.RemoveAll(dir)

	cpm, err := checkpointmanager.NewCheckpointManager(dir)
	if err != nil {
		t.Errorf("Failed to initialize checkpoint manager error=%v", err)
	}
	for _, p := range testPods {
		// Write pods should always pass unless there is an fs error
		if err := WritePod(cpm, p.pod); err != nil {
			t.Errorf("Failed to Write Pod: %v", err)
		}
	}
	// verify the correct written files are loaded from disk
	pods, err := LoadPods(cpm)
	if err != nil {
		t.Errorf("Failed to Load Pods: %v", err)
	}
	// loop through contents and check make sure
	// what was loaded matched the expected results.
	for _, p := range testPods {
		pname := p.pod.GetName()
		var lpod *v1.Pod
		for _, check := range pods {
			if check.GetName() == pname {
				lpod = check
				break
			}
		}
		if p.written {
			if lpod != nil {
				if !reflect.DeepEqual(p.pod, lpod) {
					t.Errorf("expected %#v, \ngot %#v", p.pod, lpod)
				}
			} else {
				t.Errorf("Got unexpected result for %v, should have been loaded", pname)
			}
		} else if lpod != nil {
			t.Errorf("Got unexpected result for %v, should not have been loaded", pname)
		}
		err = DeletePod(cpm, p.pod)
		if err != nil {
			t.Errorf("Failed to delete pod %v", pname)
		}
	}
	// finally validate the contents of the directory is empty.
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Errorf("Failed to read directory %v", dir)
	}
	if len(files) > 0 {
		t.Errorf("Directory %v should be empty but found %#v", dir, files)
	}
}
