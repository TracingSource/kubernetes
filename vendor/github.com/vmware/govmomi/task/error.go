package task

import "github.com/vmware/govmomi/vim25/types"

type Error struct {
	*types.LocalizedMethodFault
}

// Error returns the task's localized fault message.
func (e Error) Error() string {
	return e.LocalizedMethodFault.LocalizedMessage
}

func (e Error) Fault() types.BaseMethodFault {
	return e.LocalizedMethodFault.Fault
}
