package checkpointmanager

import (
	"fmt"
	"sync"

	"k8s.io/kubernetes/pkg/kubelet/checkpointmanager/errors"
	utilstore "k8s.io/kubernetes/pkg/kubelet/util/store"
	utilfs "k8s.io/kubernetes/pkg/util/filesystem"
)

// Checkpoint provides the process checkpoint data
type Checkpoint interface {
	MarshalCheckpoint() ([]byte, error)
	UnmarshalCheckpoint(blob []byte) error
	VerifyChecksum() error
}

// 	@implementBy: impl{}
//
// CheckpointManager provides the interface to manage checkpoint
type CheckpointManager interface {
	// CreateCheckpoint persists checkpoint in CheckpointStore.
	// checkpointKey is the key for utilstore to locate checkpoint.
	// For file backed utilstore, checkpointKey is the file name to write the checkpoint data.
	CreateCheckpoint(checkpointKey string, checkpoint Checkpoint) error
	// GetCheckpoint retrieves checkpoint from CheckpointStore.
	GetCheckpoint(checkpointKey string, checkpoint Checkpoint) error
	// WARNING: RemoveCheckpoint will not return error if checkpoint does not exist.
	RemoveCheckpoint(checkpointKey string) error
	// ListCheckpoint returns the list of existing checkpoints.
	ListCheckpoints() ([]string, error)
}

// impl is an implementation of CheckpointManager.
// It persists checkpoint in CheckpointStore
type impl struct {
	path  string
	store utilstore.Store
	mutex sync.Mutex
}

// NewCheckpointManager returns a new instance of a checkpoint manager
func NewCheckpointManager(checkpointDir string) (CheckpointManager, error) {
	fstore, err := utilstore.NewFileStore(checkpointDir, utilfs.DefaultFs{})
	if err != nil {
		return nil, err
	}

	return &impl{path: checkpointDir, store: fstore}, nil
}

// CreateCheckpoint persists checkpoint in CheckpointStore.
func (manager *impl) CreateCheckpoint(checkpointKey string, checkpoint Checkpoint) error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	blob, err := checkpoint.MarshalCheckpoint()
	if err != nil {
		return err
	}
	return manager.store.Write(checkpointKey, blob)
}

// GetCheckpoint retrieves checkpoint from CheckpointStore.
func (manager *impl) GetCheckpoint(checkpointKey string, checkpoint Checkpoint) error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	blob, err := manager.store.Read(checkpointKey)
	if err != nil {
		if err == utilstore.ErrKeyNotFound {
			return errors.ErrCheckpointNotFound
		}
		return err
	}
	err = checkpoint.UnmarshalCheckpoint(blob)
	if err == nil {
		err = checkpoint.VerifyChecksum()
	}
	return err
}

// RemoveCheckpoint will not return error if checkpoint does not exist.
func (manager *impl) RemoveCheckpoint(checkpointKey string) error {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	return manager.store.Delete(checkpointKey)
}

// ListCheckpoints returns the list of existing checkpoints.
func (manager *impl) ListCheckpoints() ([]string, error) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	keys, err := manager.store.List()
	if err != nil {
		return []string{}, fmt.Errorf("failed to list checkpoint store: %v", err)
	}
	return keys, nil
}
