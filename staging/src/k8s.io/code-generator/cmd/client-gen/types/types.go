package types

type Version string

func (v Version) String() string {
	return string(v)
}

func (v Version) NonEmpty() string {
	if v == "" {
		return "internalVersion"
	}
	return v.String()
}

type Group string

func (g Group) String() string {
	return string(g)
}

func (g Group) NonEmpty() string {
	if g == "api" {
		return "core"
	}
	return string(g)
}

type PackageVersion struct {
	Version
	// The fully qualified package, e.g. k8s.io/kubernetes/pkg/apis/apps, where the types.go is found.
	Package string
}

type GroupVersion struct {
	Group   Group
	Version Version
}

type GroupVersions struct {
	// The name of the package for this group, e.g. apps.
	PackageName string
	Group       Group
	Versions    []PackageVersion
}

// GroupVersionInfo contains all the info around a group version.
type GroupVersionInfo struct {
	Group                Group
	Version              Version
	PackageAlias         string
	GroupGoName          string
	LowerCaseGroupGoName string
}

type GroupInstallPackage struct {
	Group               Group
	InstallPackageAlias string
}
