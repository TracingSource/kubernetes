package nodeidentifier

import (
	"strings"

	"k8s.io/apiserver/pkg/authentication/user"
)

// NewDefaultNodeIdentifier returns a default NodeIdentifier implementation,
// which returns isNode=true if the user groups contain the system:nodes group
// and the user name matches the format system:node:<nodeName>, and populates
// nodeName if isNode is true
func NewDefaultNodeIdentifier() NodeIdentifier {
	return defaultNodeIdentifier{}
}

// defaultNodeIdentifier implements NodeIdentifier
type defaultNodeIdentifier struct{}

// nodeUserNamePrefix is the prefix for usernames in the form `system:node:<nodeName>`
const nodeUserNamePrefix = "system:node:"

// NodeIdentity returns isNode=true if the user groups contain the system:nodes
// group and the user name matches the format system:node:<nodeName>, and
// populates nodeName if isNode is true
func (defaultNodeIdentifier) NodeIdentity(u user.Info) (string, bool) {
	// Make sure we're a node, and can parse the node name
	if u == nil {
		return "", false
	}

	userName := u.GetName()
	if !strings.HasPrefix(userName, nodeUserNamePrefix) {
		return "", false
	}

	isNode := false
	for _, g := range u.GetGroups() {
		if g == user.NodesGroup {
			isNode = true
			break
		}
	}
	if !isNode {
		return "", false
	}

	nodeName := strings.TrimPrefix(userName, nodeUserNamePrefix)
	return nodeName, true
}
