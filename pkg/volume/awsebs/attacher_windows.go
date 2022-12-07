// +build !providerless
// +build windows

package awsebs

import (
	"fmt"
	"regexp"
	"strings"
)

// ebsnvme-id is present on AWS-provided Windows Server AMIs
// https://docs.aws.amazon.com/AWSEC2/latest/WindowsGuide/nvme-ebs-volumes.html#identify-nvme-ebs-device
const ebsnvmeID = `C:\ProgramData\Amazon\Tools\ebsnvme-id.exe`

func (attacher *awsElasticBlockStoreAttacher) getDevicePath(volumeID, partition, devicePath string) (string, error) {
	return attacher.getDiskNumber(volumeID)
}

// getDiskNumber gets the Windows disk number for a given volume ID. The disk number is needed for mounting.
// TODO This only works for Nitro-based instances
// TODO fallback to Get-Disk
func (attacher *awsElasticBlockStoreAttacher) getDiskNumber(volumeID string) (string, error) {
	// Split the ID from zone: aws://us-west-2b/vol-06d0909eb358b05f9
	split := strings.Split(volumeID, "/")
	volumeID = split[len(split)-1]

	exec := attacher.host.GetExec(awsElasticBlockStorePluginName)
	output, err := exec.Command(ebsnvmeID).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error calling ebsnvme-id.exe: %v", err)
	}
	// ebsnvme-id.exe will output a list of disks in this format:
	// ```
	// Disk Number: 1
	// Volume ID: vol-06d0909eb358b05f9
	// Device Name: /dev/xvdch
	// ```
	// Don't try to match devicePath against "Device Name" not only because volume ID is sufficient,
	// but because devicePath may change between Linux & Windows formats between WaitForAttach calls.
	// The first attach and mount, WaitForAttach gets devicePath as the Linux format /dev/xvdch. Then
	// WaitForAttach returns the disk number as the "right" devicePath and that is persisted to ASW.
	// In subsequent mounts of the same disk, WaitForAttach gets devicePath as the Windows format it
	// returned the first time.
	diskRe := regexp.MustCompile(
		`Disk Number: (\d+)\s*` +
			`Volume ID: ` + volumeID + `\s*`)
	matches := diskRe.FindStringSubmatch(string(output))
	if len(matches) != 2 {
		return "", fmt.Errorf("disk not found in ebsnvme-id.exe output: %q", string(output))
	}
	return matches[1], nil
}
