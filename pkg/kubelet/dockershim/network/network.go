package network

// TODO: Consider making this value configurable.
const DefaultInterfaceName = "eth0"

// CNITimeoutSec is set to be slightly less than 240sec/4mins, which is the default remote runtime request timeout.
const CNITimeoutSec = 220

// UseDefaultMTU is a marker value that indicates the plugin should determine its own MTU
// It is the zero value, so a non-initialized value will mean "UseDefault"
const UseDefaultMTU = 0
