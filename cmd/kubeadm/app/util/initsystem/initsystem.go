package initsystem

// InitSystem is the interface that describe behaviors of an init system
type InitSystem interface {
	// return a string describing how to enable a service
	EnableCommand(service string) string

	// ServiceStart tries to start a specific service
	ServiceStart(service string) error

	// ServiceStop tries to stop a specific service
	ServiceStop(service string) error

	// ServiceRestart tries to reload the environment and restart the specific service
	ServiceRestart(service string) error

	// ServiceExists ensures the service is defined for this init system.
	ServiceExists(service string) bool

	// ServiceIsEnabled ensures the service is enabled to start on each boot.
	ServiceIsEnabled(service string) bool

	// ServiceIsActive ensures the service is running, or attempting to run. (crash looping in the case of kubelet)
	ServiceIsActive(service string) bool
}
