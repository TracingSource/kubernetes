package node

// intSet maintains a set of ints, and supports promoting and culling the previous generation.
// this allows tracking a large, mostly-stable set without constantly reallocating the entire set.
type intSet struct {
	currentGeneration byte
	members           map[int]byte
}

func newIntSet() *intSet {
	return &intSet{members: map[int]byte{}}
}

// has returns true if the specified int is in the set.
// it is safe to call concurrently, but must not be called concurrently with any of the other methods.
func (s *intSet) has(i int) bool {
	if s == nil {
		return false
	}
	_, present := s.members[i]
	return present
}

// startNewGeneration begins a new generation.
// it must be followed by a call to mark() for every member of the generation,
// then a call to sweep() to remove members not present in the generation.
// it is not thread-safe.
func (s *intSet) startNewGeneration() {
	s.currentGeneration++
}

// mark indicates the specified int belongs to the current generation.
// it is not thread-safe.
func (s *intSet) mark(i int) {
	s.members[i] = s.currentGeneration
}

// sweep removes items not in the current generation.
// it is not thread-safe.
func (s *intSet) sweep() {
	for k, v := range s.members {
		if v != s.currentGeneration {
			delete(s.members, k)
		}
	}
}
