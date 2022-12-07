// +build !providerless

package azure

import (
	"sync"
)

// lockMap used to lock on entries
type lockMap struct {
	sync.Mutex
	mutexMap map[string]*sync.Mutex
}

// NewLockMap returns a new lock map
func newLockMap() *lockMap {
	return &lockMap{
		mutexMap: make(map[string]*sync.Mutex),
	}
}

// LockEntry acquires a lock associated with the specific entry
func (lm *lockMap) LockEntry(entry string) {
	lm.Lock()
	// check if entry does not exists, then add entry
	if _, exists := lm.mutexMap[entry]; !exists {
		lm.addEntry(entry)
	}

	lm.Unlock()
	lm.lockEntry(entry)
}

// UnlockEntry release the lock associated with the specific entry
func (lm *lockMap) UnlockEntry(entry string) {
	lm.Lock()
	defer lm.Unlock()

	if _, exists := lm.mutexMap[entry]; !exists {
		return
	}
	lm.unlockEntry(entry)
}

func (lm *lockMap) addEntry(entry string) {
	lm.mutexMap[entry] = &sync.Mutex{}
}

func (lm *lockMap) lockEntry(entry string) {
	lm.mutexMap[entry].Lock()
}

func (lm *lockMap) unlockEntry(entry string) {
	lm.mutexMap[entry].Unlock()
}
