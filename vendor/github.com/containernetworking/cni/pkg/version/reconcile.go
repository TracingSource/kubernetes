package version

import "fmt"

type ErrorIncompatible struct {
	Config    string
	Supported []string
}

func (e *ErrorIncompatible) Details() string {
	return fmt.Sprintf("config is %q, plugin supports %q", e.Config, e.Supported)
}

func (e *ErrorIncompatible) Error() string {
	return fmt.Sprintf("incompatible CNI versions: %s", e.Details())
}

type Reconciler struct{}

func (r *Reconciler) Check(configVersion string, pluginInfo PluginInfo) *ErrorIncompatible {
	return r.CheckRaw(configVersion, pluginInfo.SupportedVersions())
}

func (*Reconciler) CheckRaw(configVersion string, supportedVersions []string) *ErrorIncompatible {
	for _, supportedVersion := range supportedVersions {
		if configVersion == supportedVersion {
			return nil
		}
	}

	return &ErrorIncompatible{
		Config:    configVersion,
		Supported: supportedVersions,
	}
}
