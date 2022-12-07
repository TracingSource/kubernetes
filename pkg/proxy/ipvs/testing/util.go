package testing

import (
	utilipset "k8s.io/kubernetes/pkg/util/ipset"
)

// ExpectedVirtualServer is the expected ipvs rules with VirtualServer and RealServer
// VSNum is the expected ipvs virtual server number
// IP:Port protocol is the expected ipvs vs info
// RS is the RealServer of this expected VirtualServer
type ExpectedVirtualServer struct {
	VSNum    int
	IP       string
	Port     uint16
	Protocol string
	RS       []ExpectedRealServer
}

// ExpectedRealServer is the expected ipvs RealServer
type ExpectedRealServer struct {
	IP   string
	Port uint16
}

// ExpectedIptablesChain is a map of expected iptables chain and jump rules
type ExpectedIptablesChain map[string][]ExpectedIptablesRule

// ExpectedIptablesRule is the expected iptables rules with jump chain and match ipset name
type ExpectedIptablesRule struct {
	JumpChain string
	MatchSet  string
}

// ExpectedIPSet is the expected ipset with set name and entries name
type ExpectedIPSet map[string][]*utilipset.Entry
