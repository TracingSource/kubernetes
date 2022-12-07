package rangeallocation

import (
	api "k8s.io/kubernetes/pkg/apis/core"
)

// RangeRegistry is a registry that can retrieve or persist a RangeAllocation object.
type RangeRegistry interface {
	// Get returns the latest allocation, an empty object if no allocation has been made,
	// or an error if the allocation could not be retrieved.
	Get() (*api.RangeAllocation, error)
	// CreateOrUpdate should create or update the provide allocation, unless a conflict
	// has occurred since the item was last created.
	CreateOrUpdate(*api.RangeAllocation) error
}
