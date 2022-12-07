package wardle

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FlunderList is a list of Flunder objects.
type FlunderList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Flunder
}

type ReferenceType string

const (
	FlunderReferenceType = ReferenceType("Flunder")
	FischerReferenceType = ReferenceType("Fischer")
)

type FlunderSpec struct {
	// A name of another flunder, mutually exclusive to the FischerReference.
	FlunderReference string
	// A name of a fischer, mutually exclusive to the FlunderReference.
	FischerReference string
	// The reference type.
	ReferenceType ReferenceType
}

type FlunderStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Flunder struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   FlunderSpec
	Status FlunderStatus
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Fischer struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	// DisallowedFlunders holds a list of Flunder.Names that are disallowed.
	DisallowedFlunders []string
}

// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// FischerList is a list of Fischer objects.
type FischerList struct {
	metav1.TypeMeta
	metav1.ListMeta

	// Items is a list of Fischers
	Items []Fischer
}
