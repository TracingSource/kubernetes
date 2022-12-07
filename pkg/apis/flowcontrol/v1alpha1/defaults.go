package v1alpha1

import (
	"k8s.io/api/flowcontrol/v1alpha1"
)

// Default settings for flow-schema
const (
	FlowSchemaDefaultMatchingPrecedence int32 = 1000
)

// Default settings for priority-level-configuration
const (
	PriorityLevelConfigurationDefaultHandSize                 int32 = 8
	PriorityLevelConfigurationDefaultQueues                   int32 = 64
	PriorityLevelConfigurationDefaultQueueLengthLimit         int32 = 50
	PriorityLevelConfigurationDefaultAssuredConcurrencyShares int32 = 30
)

// SetDefaults_FlowSchema sets default values for flow schema
func SetDefaults_FlowSchemaSpec(spec *v1alpha1.FlowSchemaSpec) {
	if spec.MatchingPrecedence == 0 {
		spec.MatchingPrecedence = FlowSchemaDefaultMatchingPrecedence
	}
}

func SetDefaults_LimitedPriorityLevelConfiguration(lplc *v1alpha1.LimitedPriorityLevelConfiguration) {
	if lplc.AssuredConcurrencyShares == 0 {
		lplc.AssuredConcurrencyShares = PriorityLevelConfigurationDefaultAssuredConcurrencyShares
	}
}

// SetDefaults_FlowSchema sets default values for flow schema
func SetDefaults_QueuingConfiguration(cfg *v1alpha1.QueuingConfiguration) {
	if cfg.HandSize == 0 {
		cfg.HandSize = PriorityLevelConfigurationDefaultHandSize
	}
	if cfg.Queues == 0 {
		cfg.Queues = PriorityLevelConfigurationDefaultQueues
	}
	if cfg.QueueLengthLimit == 0 {
		cfg.QueueLengthLimit = PriorityLevelConfigurationDefaultQueueLengthLimit
	}
}
