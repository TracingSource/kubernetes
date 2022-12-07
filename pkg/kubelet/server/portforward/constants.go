// Package portforward contains server-side logic for handling port forwarding requests.
package portforward

// ProtocolV1Name is the name of the subprotocol used for port forwarding.
const ProtocolV1Name = "portforward.k8s.io"

// SupportedProtocols are the supported port forwarding protocols.
var SupportedProtocols = []string{ProtocolV1Name}
