/*
Copyright (c) 2016 VMware, Inc. All Rights Reserved.


*/

package mo

// Entity is the interface that is implemented by all managed objects
// that extend ManagedEntity.
type Entity interface {
	Reference
	Entity() *ManagedEntity
}
