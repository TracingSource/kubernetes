// Package types defines types used only by volume components
package types

import (
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/runtime"
)

// UniquePodName defines the type to key pods off of
type UniquePodName types.UID

// UniquePVCName defines the type to key pvc off
type UniquePVCName types.UID

// GeneratedOperations contains the operation that is created as well as
// supporting functions required for the operation executor
type GeneratedOperations struct {
	// Name of operation - could be used for resetting shared exponential backoff
	OperationName     string
	// caller: 在 GeneratedOperations{} 对象中, 通过其 Run() 方法调用.
	OperationFunc     func() (eventErr error, detailedErr error)
	EventRecorderFunc func(*error)
	CompleteFunc      func(*error)
}

// Run ...
// caller: 
// 	1. pkg/volume/util/nestedpendingoperations/nestedpendingoperations.go -> nestedPendingOperations.Run()
//  在其中的 go 协程中被调用.
//
// Run executes the operations and its supporting functions
func (o *GeneratedOperations) Run() (eventErr, detailedErr error) {
	if o.CompleteFunc != nil {
		defer o.CompleteFunc(&detailedErr)
	}
	if o.EventRecorderFunc != nil {
		defer o.EventRecorderFunc(&eventErr)
	}
	// Handle panic, if any, from operationFunc()
	defer runtime.RecoverFromPanic(&detailedErr)
	return o.OperationFunc()
}

const (
	// VolumeResizerKey is key that will be used to store resizer used
	// for resizing PVC. The generated key/value pair will be added
	// as a annotation to the PVC.
	VolumeResizerKey = "volume.kubernetes.io/storage-resizer"
)
