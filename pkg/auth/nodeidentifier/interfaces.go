package nodeidentifier

import (
	"k8s.io/apiserver/pkg/authentication/user"
)

// NodeIdentifier determines node information from a given user
type NodeIdentifier interface {
	// NodeIdentity determines node information from the given user.Info.
	// nodeName is the name of the Node API object associated with the user.Info,
	// and may be empty if a specific node cannot be determined.
	// isNode is true if the user.Info represents an identity issued to a node.
	NodeIdentity(user.Info) (nodeName string, isNode bool)
}
