package core

import "encoding/json"

// This file implements json marshaling/unmarshaling interfaces on objects that are currently marshaled into annotations
// to prevent anyone from marshaling these internal structs.

var _ = json.Marshaler(&AvoidPods{})
var _ = json.Unmarshaler(&AvoidPods{})

// MarshalJSON panics to prevent marshalling of internal structs
func (AvoidPods) MarshalJSON() ([]byte, error) { panic("do not marshal internal struct") }

// UnmarshalJSON panics to prevent unmarshalling of internal structs
func (*AvoidPods) UnmarshalJSON([]byte) error { panic("do not unmarshal to internal struct") }
