// +build !providerless

package azure

import (
	"testing"
	"time"
)

func TestSimpleLockEntry(t *testing.T) {
	testLockMap := newLockMap()

	callbackChan1 := make(chan interface{})
	go testLockMap.lockAndCallback(t, "entry1", callbackChan1)
	ensureCallbackHappens(t, callbackChan1)
}

func TestSimpleLockUnlockEntry(t *testing.T) {
	testLockMap := newLockMap()

	callbackChan1 := make(chan interface{})
	go testLockMap.lockAndCallback(t, "entry1", callbackChan1)
	ensureCallbackHappens(t, callbackChan1)
	testLockMap.UnlockEntry("entry1")
}

func TestConcurrentLockEntry(t *testing.T) {
	testLockMap := newLockMap()

	callbackChan1 := make(chan interface{})
	callbackChan2 := make(chan interface{})

	go testLockMap.lockAndCallback(t, "entry1", callbackChan1)
	ensureCallbackHappens(t, callbackChan1)

	go testLockMap.lockAndCallback(t, "entry1", callbackChan2)
	ensureNoCallback(t, callbackChan2)

	testLockMap.UnlockEntry("entry1")
	ensureCallbackHappens(t, callbackChan2)
	testLockMap.UnlockEntry("entry1")
}

func (lm *lockMap) lockAndCallback(t *testing.T, entry string, callbackChan chan<- interface{}) {
	lm.LockEntry(entry)
	callbackChan <- true
}

var callbackTimeout = 2 * time.Second

func ensureCallbackHappens(t *testing.T, callbackChan <-chan interface{}) bool {
	select {
	case <-callbackChan:
		return true
	case <-time.After(callbackTimeout):
		t.Fatalf("timed out waiting for callback")
		return false
	}
}

func ensureNoCallback(t *testing.T, callbackChan <-chan interface{}) bool {
	select {
	case <-callbackChan:
		t.Fatalf("unexpected callback")
		return false
	case <-time.After(callbackTimeout):
		return true
	}
}
