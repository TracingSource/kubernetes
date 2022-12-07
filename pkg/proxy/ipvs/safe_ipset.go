package ipvs

import (
	"sync"

	"k8s.io/kubernetes/pkg/util/ipset"
)

type safeIpset struct {
	ipset ipset.Interface
	mu    sync.Mutex
}

func newSafeIpset(ipset ipset.Interface) ipset.Interface {
	return &safeIpset{
		ipset: ipset,
	}
}

// FlushSet deletes all entries from a named set.
func (s *safeIpset) FlushSet(set string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ipset.FlushSet(set)
}

// DestroySet deletes a named set.
func (s *safeIpset) DestroySet(set string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ipset.DestroySet(set)
}

// DestroyAllSets deletes all sets.
func (s *safeIpset) DestroyAllSets() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ipset.DestroyAllSets()
}

// CreateSet creates a new set.  It will ignore error when the set already exists if ignoreExistErr=true.
func (s *safeIpset) CreateSet(set *ipset.IPSet, ignoreExistErr bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ipset.CreateSet(set, ignoreExistErr)
}

// AddEntry adds a new entry to the named set.  It will ignore error when the entry already exists if ignoreExistErr=true.
func (s *safeIpset) AddEntry(entry string, set *ipset.IPSet, ignoreExistErr bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ipset.AddEntry(entry, set, ignoreExistErr)
}

// DelEntry deletes one entry from the named set
func (s *safeIpset) DelEntry(entry string, set string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ipset.DelEntry(entry, set)
}

// Test test if an entry exists in the named set
func (s *safeIpset) TestEntry(entry string, set string) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ipset.TestEntry(entry, set)
}

// ListEntries lists all the entries from a named set
func (s *safeIpset) ListEntries(set string) ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ipset.ListEntries(set)
}

// ListSets list all set names from kernel
func (s *safeIpset) ListSets() ([]string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ipset.ListSets()
}

// GetVersion returns the "X.Y" version string for ipset.
func (s *safeIpset) GetVersion() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.ipset.GetVersion()
}
