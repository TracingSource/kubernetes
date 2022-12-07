// +build linux

package v1beta1

// DefaultEvictionHard includes default options for hard eviction.
var DefaultEvictionHard = map[string]string{
	"memory.available":  "100Mi",
	"nodefs.available":  "10%",
	"nodefs.inodesFree": "5%",
	"imagefs.available": "15%",
}
