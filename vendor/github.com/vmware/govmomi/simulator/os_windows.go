package simulator

import "os"

func (ds *Datastore) stat() error {
	info := ds.Info.GetDatastoreInfo()

	_, err := os.Stat(info.Url)
	return err
}
