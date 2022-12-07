// +build windows

package winstats

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
)

//OSInfo is a convenience class for retrieving Windows OS information
type OSInfo struct {
	BuildNumber, ProductName        string
	MajorVersion, MinorVersion, UBR uint64
}

// GetOSInfo reads Windows version information from the registry
func GetOSInfo() (*OSInfo, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, `SOFTWARE\Microsoft\Windows NT\CurrentVersion`, registry.QUERY_VALUE)
	if err != nil {
		return nil, err
	}
	defer k.Close()

	buildNumber, _, err := k.GetStringValue("CurrentBuildNumber")
	if err != nil {
		return nil, err
	}

	majorVersionNumber, _, err := k.GetIntegerValue("CurrentMajorVersionNumber")
	if err != nil {
		return nil, err
	}

	minorVersionNumber, _, err := k.GetIntegerValue("CurrentMinorVersionNumber")
	if err != nil {
		return nil, err
	}

	revision, _, err := k.GetIntegerValue("UBR")
	if err != nil {
		return nil, err
	}

	productName, _, err := k.GetStringValue("ProductName")
	if err != nil {
		return nil, nil
	}

	return &OSInfo{
		BuildNumber:  buildNumber,
		ProductName:  productName,
		MajorVersion: majorVersionNumber,
		MinorVersion: minorVersionNumber,
		UBR:          revision,
	}, nil
}

//GetPatchVersion returns full OS version with patch
func (o *OSInfo) GetPatchVersion() string {
	return fmt.Sprintf("%d.%d.%s.%d", o.MajorVersion, o.MinorVersion, o.BuildNumber, o.UBR)
}

//GetBuild returns OS version upto build number
func (o *OSInfo) GetBuild() string {
	return fmt.Sprintf("%d.%d.%s", o.MajorVersion, o.MinorVersion, o.BuildNumber)
}
