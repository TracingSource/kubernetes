package utils

import (
	"io/ioutil"

	"k8s.io/klog"
)

func MakeTempDirOrDie(prefix string, baseDir string) string {
	if baseDir == "" {
		baseDir = "/tmp"
	}
	tempDir, err := ioutil.TempDir(baseDir, prefix)
	if err != nil {
		klog.Fatalf("Can't make a temp rootdir: %v", err)
	}
	return tempDir
}
