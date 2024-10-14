package cache

// 	@implementBy: pkg/kubelet/cm/devicemanager/manager.go -> ManagerImpl{}
//
// PluginHandler is an interface a client of the pluginwatcher API needs to implement in
// order to consume plugins
// The PluginHandler follows the simple following state machine:
//
//                         +--------------------------------------+
//                         |            ReRegistration            |
//                         | Socket created with same plugin name |
//                         |                                      |
//                         |                                      |
//    Socket Created       v                                      +        Socket Deleted
// +------------------> Validate +---------------------------> Register +------------------> DeRegister
//                         +                                      +                              +
//                         |                                      |                              |
//                         | Error                                | Error                        |
//                         |                                      |                              |
//                         v                                      v                              v
//                        Out                                    Out                            Out
//
// The pluginwatcher module follows strictly and sequentially this state machine for each *plugin name*.
// e.g: If you are Registering a plugin foo, you cannot get a DeRegister call for plugin foo
//      until the Register("foo") call returns. Nor will you get a Validate("foo", "Different endpoint", ...)
//      call until the Register("foo") call returns.
//
// ReRegistration: Socket created with same plugin name, usually for a plugin update
// e.g: plugin with name foo registers at foo.com/foo-1.9.7 later a plugin with name foo
//      registers at foo.com/foo-1.9.9
//
// DeRegistration: When ReRegistration happens only the deletion of the new socket will trigger a DeRegister call
type PluginHandler interface {
	// Validate returns an error if the information provided by
	// the potential plugin is erroneous (unsupported version, ...)
	ValidatePlugin(pluginName string, endpoint string, versions []string) error
	// RegisterPlugin is called so that the plugin can be register by any
	// plugin consumer
	// Error encountered here can still be Notified to the plugin.
	RegisterPlugin(pluginName, endpoint string, versions []string) error
	// DeRegister is called once the pluginwatcher observes that the socket has
	// been deleted.
	DeRegisterPlugin(pluginName string)
}
