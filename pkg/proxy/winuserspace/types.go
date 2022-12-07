package winuserspace

import (
	"fmt"

	"k8s.io/apimachinery/pkg/types"
)

// ServicePortPortalName carries a namespace + name + portname + portalip.  This is the unique
// identifier for a windows service port portal.
type ServicePortPortalName struct {
	types.NamespacedName
	Port         string
	PortalIPName string
}

func (spn ServicePortPortalName) String() string {
	return fmt.Sprintf("%s:%s:%s", spn.NamespacedName.String(), spn.Port, spn.PortalIPName)
}
