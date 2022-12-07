package checksum

import (
	"hash/fnv"

	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager/errors"
	hashutil "k8s.io/kubernetes/pkg/util/hash"
)

// Checksum is the data to be stored as checkpoint
type Checksum uint64

// Verify verifies that passed checksum is same as calculated checksum
func (cs Checksum) Verify(data interface{}) error {
	if cs != New(data) {
		return errors.ErrCorruptCheckpoint
	}
	return nil
}

// New returns the Checksum of checkpoint data
func New(data interface{}) Checksum {
	return Checksum(getChecksum(data))
}

// Get returns calculated checksum of checkpoint data
func getChecksum(data interface{}) uint64 {
	hash := fnv.New32a()
	hashutil.DeepHashObject(hash, data)
	return uint64(hash.Sum32())
}
