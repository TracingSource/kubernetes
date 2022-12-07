package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FlunderList is a list of Flunder objects.
type FlunderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []Flunder `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type ReferenceType string

const (
	FlunderReferenceType = ReferenceType("Flunder")
	FischerReferenceType = ReferenceType("Fischer")
)

type FlunderSpec struct {
	// A name of another flunder or fischer, depending on the reference type.
	Reference string `json:"reference,omitempty" protobuf:"bytes,1,opt,name=reference"`
	// The reference type, defaults to "Flunder" if reference is set.
	ReferenceType *ReferenceType `json:"referenceType,omitempty" protobuf:"bytes,2,opt,name=referenceType"`
}

type FlunderStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Flunder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   FlunderSpec   `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status FlunderStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Fischer struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// DisallowedFlunders holds a list of Flunder.Names that are disallowed.
	DisallowedFlunders []string `json:"disallowedFlunders,omitempty" protobuf:"bytes,2,rep,name=disallowedFlunders"`
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FischerList is a list of Fischer objects.
type FischerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []Fischer `json:"items" protobuf:"bytes,2,rep,name=items"`
}
