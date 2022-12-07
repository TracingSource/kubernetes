// +build windows

package backend

import bolt "go.etcd.io/bbolt"

var boltOpenOptions *bolt.Options = nil

// setting mmap size != 0 on windows will allocate the entire
// mmap size for the file, instead of growing it. So, force 0.

func (bcfg *BackendConfig) mmapSize() int { return 0 }
