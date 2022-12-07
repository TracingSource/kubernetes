package core

import (
	"k8s.io/apimachinery/pkg/api/resource"
)

func (rn ResourceName) String() string {
	return string(rn)
}

// CPU returns the CPU limit if specified.
func (rl *ResourceList) CPU() *resource.Quantity {
	if val, ok := (*rl)[ResourceCPU]; ok {
		return &val
	}
	return &resource.Quantity{Format: resource.DecimalSI}
}

// Memory returns the Memory limit if specified.
func (rl *ResourceList) Memory() *resource.Quantity {
	if val, ok := (*rl)[ResourceMemory]; ok {
		return &val
	}
	return &resource.Quantity{Format: resource.BinarySI}
}

// Pods returns the list of pods
func (rl *ResourceList) Pods() *resource.Quantity {
	if val, ok := (*rl)[ResourcePods]; ok {
		return &val
	}
	return &resource.Quantity{}
}

// StorageEphemeral returns the list of ephemeral storage volumes, if any
func (rl *ResourceList) StorageEphemeral() *resource.Quantity {
	if val, ok := (*rl)[ResourceEphemeralStorage]; ok {
		return &val
	}
	return &resource.Quantity{}
}
