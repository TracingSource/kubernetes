package util

import (
	"sort"

	"k8s.io/kubernetes/pkg/apis/flowcontrol"
)

var _ sort.Interface = FlowSchemaSequence{}

// FlowSchemaSequence holds sorted set of pointers to FlowSchema objects.
// FlowSchemaSequence implements `sort.Interface`
type FlowSchemaSequence []*flowcontrol.FlowSchema

func (s FlowSchemaSequence) Len() int {
	return len(s)
}

func (s FlowSchemaSequence) Less(i, j int) bool {
	// the flow-schema w/ lower matching-precedence is prior
	if ip, jp := s[i].Spec.MatchingPrecedence, s[j].Spec.MatchingPrecedence; ip != jp {
		return ip < jp
	}
	// sort alphabetically
	return s[i].Name < s[j].Name
}

func (s FlowSchemaSequence) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
