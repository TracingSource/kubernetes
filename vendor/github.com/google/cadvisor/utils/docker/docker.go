
package docker

import (
	"fmt"
	"os"
	"strings"

	dockertypes "github.com/docker/docker/api/types"
)

const (
	DockerInfoDriver          = "Driver"
	DockerInfoDriverStatus    = "DriverStatus"
	DriverStatusPoolName      = "Pool Name"
	DriverStatusDataLoopFile  = "Data loop file"
	DriverStatusMetadataFile  = "Metadata file"
	DriverStatusParentDataset = "Parent Dataset"
)

func DriverStatusValue(status [][2]string, target string) string {
	for _, v := range status {
		if strings.EqualFold(v[0], target) {
			return v[1]
		}
	}

	return ""
}

func DockerThinPoolName(info dockertypes.Info) (string, error) {
	poolName := DriverStatusValue(info.DriverStatus, DriverStatusPoolName)
	if len(poolName) == 0 {
		return "", fmt.Errorf("Could not get devicemapper pool name")
	}

	return poolName, nil
}

func DockerMetadataDevice(info dockertypes.Info) (string, error) {
	metadataDevice := DriverStatusValue(info.DriverStatus, DriverStatusMetadataFile)
	if len(metadataDevice) != 0 {
		return metadataDevice, nil
	}

	poolName, err := DockerThinPoolName(info)
	if err != nil {
		return "", err
	}

	metadataDevice = fmt.Sprintf("/dev/mapper/%s_tmeta", poolName)

	if _, err := os.Stat(metadataDevice); err != nil {
		return "", err
	}

	return metadataDevice, nil
}

func DockerZfsFilesystem(info dockertypes.Info) (string, error) {
	filesystem := DriverStatusValue(info.DriverStatus, DriverStatusParentDataset)
	if len(filesystem) == 0 {
		return "", fmt.Errorf("Could not get zfs filesystem")
	}

	return filesystem, nil
}
