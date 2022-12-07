package quorum

import (
	"math"
	"strconv"
)

// Index is a Raft log position.
type Index uint64

func (i Index) String() string {
	if i == math.MaxUint64 {
		return "∞"
	}
	return strconv.FormatUint(uint64(i), 10)
}

// AckedIndexer allows looking up a commit index for a given ID of a voter
// from a corresponding MajorityConfig.
type AckedIndexer interface {
	AckedIndex(voterID uint64) (idx Index, found bool)
}

type mapAckIndexer map[uint64]Index

func (m mapAckIndexer) AckedIndex(id uint64) (Index, bool) {
	idx, ok := m[id]
	return idx, ok
}

// VoteResult indicates the outcome of a vote.
//
//go:generate stringer -type=VoteResult
type VoteResult uint8

const (
	// VotePending indicates that the decision of the vote depends on future
	// votes, i.e. neither "yes" or "no" has reached quorum yet.
	VotePending VoteResult = 1 + iota
	// VoteLost indicates that the quorum has voted "no".
	VoteLost
	// VoteWon indicates that the quorum has voted "yes".
	VoteWon
)
