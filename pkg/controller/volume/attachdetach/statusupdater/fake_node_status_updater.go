package statusupdater

import (
	"fmt"
)

func NewFakeNodeStatusUpdater(returnError bool) NodeStatusUpdater {
	return &fakeNodeStatusUpdater{
		returnError: returnError,
	}
}

type fakeNodeStatusUpdater struct {
	returnError bool
}

func (fnsu *fakeNodeStatusUpdater) UpdateNodeStatuses() error {
	if fnsu.returnError {
		return fmt.Errorf("fake error on update node status")
	}

	return nil
}
