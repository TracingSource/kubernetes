package cache

import (
	"fmt"
	"testing"
	"time"

	"k8s.io/apimachinery/pkg/util/clock"
	expirationcache "k8s.io/client-go/tools/cache"
)

type testObject struct {
	key string
	val string
}

// A fake objectCache for unit test.
func NewFakeObjectCache(f func() (interface{}, error), ttl time.Duration, clock clock.Clock) *ObjectCache {
	ttlPolicy := &expirationcache.TTLPolicy{TTL: ttl, Clock: clock}
	deleteChan := make(chan string, 1)
	return &ObjectCache{
		updater: f,
		cache:   expirationcache.NewFakeExpirationStore(stringKeyFunc, deleteChan, ttlPolicy, clock),
	}
}

func TestAddAndGet(t *testing.T) {
	testObj := testObject{
		key: "foo",
		val: "bar",
	}
	objectCache := NewFakeObjectCache(func() (interface{}, error) {
		return nil, fmt.Errorf("Unexpected Error: updater should never be called in this test")
	}, 1*time.Hour, clock.NewFakeClock(time.Now()))

	err := objectCache.Add(testObj.key, testObj.val)
	if err != nil {
		t.Errorf("Unable to add obj %#v by key: %s", testObj, testObj.key)
	}
	value, err := objectCache.Get(testObj.key)
	if err != nil {
		t.Errorf("Unable to get obj %#v by key: %s", testObj, testObj.key)
	}
	if value.(string) != testObj.val {
		t.Errorf("Expected to get cached value: %#v, but got: %s", testObj.val, value.(string))
	}

}

func TestExpirationBasic(t *testing.T) {
	unexpectedVal := "bar"
	expectedVal := "bar2"

	testObj := testObject{
		key: "foo",
		val: unexpectedVal,
	}

	fakeClock := clock.NewFakeClock(time.Now())

	objectCache := NewFakeObjectCache(func() (interface{}, error) {
		return expectedVal, nil
	}, 1*time.Second, fakeClock)

	err := objectCache.Add(testObj.key, testObj.val)
	if err != nil {
		t.Errorf("Unable to add obj %#v by key: %s", testObj, testObj.key)
	}

	// sleep 2s so cache should be expired.
	fakeClock.Sleep(2 * time.Second)

	value, err := objectCache.Get(testObj.key)
	if err != nil {
		t.Errorf("Unable to get obj %#v by key: %s", testObj, testObj.key)
	}
	if value.(string) != expectedVal {
		t.Errorf("Expected to get cached value: %#v, but got: %s", expectedVal, value.(string))
	}
}
