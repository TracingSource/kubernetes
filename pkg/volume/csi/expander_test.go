package csi

import (
	"os"
	"testing"

	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/kubernetes/pkg/volume"
)

func TestNodeExpand(t *testing.T) {
	tests := []struct {
		name          string
		nodeExpansion bool
		nodeStageSet  bool
		volumePhase   volume.CSIVolumePhaseType
		success       bool
		fsVolume      bool
	}{
		{
			name:    "when node expansion is not set",
			success: false,
		},
		{
			name:          "when nodeExpansion=on, nodeStage=on, volumePhase=staged",
			nodeExpansion: true,
			nodeStageSet:  true,
			volumePhase:   volume.CSIVolumeStaged,
			success:       true,
			fsVolume:      true,
		},
		{
			name:          "when nodeExpansion=on, nodeStage=off, volumePhase=staged",
			nodeExpansion: true,
			volumePhase:   volume.CSIVolumeStaged,
			success:       false,
			fsVolume:      true,
		},
		{
			name:          "when nodeExpansion=on, nodeStage=on, volumePhase=published",
			nodeExpansion: true,
			nodeStageSet:  true,
			volumePhase:   volume.CSIVolumePublished,
			success:       true,
			fsVolume:      true,
		},
		{
			name:          "when nodeExpansion=on, nodeStage=off, volumePhase=published",
			nodeExpansion: true,
			volumePhase:   volume.CSIVolumePublished,
			success:       true,
			fsVolume:      true,
		},
		{
			name:          "when nodeExpansion=on, nodeStage=off, volumePhase=published, fsVolume=false",
			nodeExpansion: true,
			volumePhase:   volume.CSIVolumePublished,
			success:       true,
			fsVolume:      false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			plug, tmpDir := newTestPlugin(t, nil)
			defer os.RemoveAll(tmpDir)

			spec := volume.NewSpecFromPersistentVolume(makeTestPV("test-pv", 10, "expandable", "test-vol"), false)

			newSize, _ := resource.ParseQuantity("20Gi")

			resizeOptions := volume.NodeResizeOptions{
				VolumeSpec:      spec,
				NewSize:         newSize,
				DeviceMountPath: "/foo/bar",
				DevicePath:      "/mnt/foobar",
				CSIVolumePhase:  tc.volumePhase,
			}
			csiSource, _ := getCSISourceFromSpec(resizeOptions.VolumeSpec)

			csClient := setupClientWithExpansion(t, tc.nodeStageSet, tc.nodeExpansion)

			ok, err := plug.nodeExpandWithClient(resizeOptions, csiSource, csClient, tc.fsVolume)
			if ok != tc.success {
				if err != nil {
					t.Errorf("For %s : expected %v got %v with %v", tc.name, tc.success, ok, err)
				} else {
					t.Errorf("For %s : expected %v got %v", tc.name, tc.success, ok)
				}

			}
		})
	}
}
