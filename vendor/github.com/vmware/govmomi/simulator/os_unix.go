//+build !windows

package simulator

import "syscall"

func (ds *Datastore) stat() error {
	info := ds.Info.GetDatastoreInfo()
	var stat syscall.Statfs_t

	err := syscall.Statfs(info.Url, &stat)
	if err != nil {
		return err
	}

	info.FreeSpace = int64(stat.Bfree * uint64(stat.Bsize))

	ds.Summary.FreeSpace = info.FreeSpace
	ds.Summary.Capacity = int64(stat.Blocks * uint64(stat.Bsize))

	return nil
}
