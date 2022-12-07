package stats

import (
	"time"
)

// ResourceAnalyzer provides statistics on node resource consumption
type ResourceAnalyzer interface {
	Start()

	fsResourceAnalyzerInterface
	SummaryProvider
}

// resourceAnalyzer implements ResourceAnalyzer
type resourceAnalyzer struct {
	*fsResourceAnalyzer
	SummaryProvider
}

var _ ResourceAnalyzer = &resourceAnalyzer{}

// NewResourceAnalyzer 这里创建了 fsResourceAnalyzer 对象, 
// 且返回的结果被赋值给了 Kubelet.resourceAnalyzer 成员.
//
// NewResourceAnalyzer returns a new ResourceAnalyzer
func NewResourceAnalyzer(
	statsProvider Provider, calVolumeFrequency time.Duration,
) ResourceAnalyzer {
	fsAnalyzer := newFsResourceAnalyzer(statsProvider, calVolumeFrequency)
	summaryProvider := NewSummaryProvider(statsProvider)
	return &resourceAnalyzer{fsAnalyzer, summaryProvider}
}

// caller: 
// 	1. pkg/kubelet/kubelet__init.go -> initializeModules()
//
// Start starts background functions necessary for the ResourceAnalyzer to function
func (ra *resourceAnalyzer) Start() {
	ra.fsResourceAnalyzer.Start()
}
