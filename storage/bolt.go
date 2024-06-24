package storage

import (
	"sync"

	"github.com/boltdb/bolt"
)

const (
	ExtensionBucket = "ExtensionBucket"
)

var (
	boltOnce     sync.Once
	boltInstance *database = nil
)

type database struct {
	bolt *bolt.DB
}

func DB() *database {
	if boltInstance == nil {
		boltOnce.Do(func() {
			db, err := bolt.Open("pair.db", 0600, nil)
			if err != nil {
				panic("do not open bolt database")
			}

			boltInstance = &database{
				bolt: db,
			}

			err = boltInstance.bolt.Update(func(tx *bolt.Tx) error {

				_, err := tx.CreateBucketIfNotExists([]byte(ExtensionBucket))
				return err
			})
			if err != nil {
				panic("do not create bolt bucket")
			}

		})
	}

	return boltInstance
}

func (d *database) Close() error {
	return boltInstance.bolt.Close()
}

func (d *database) SaveExtendsion(uuid string, extension string) error {
	return d.bolt.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ExtensionBucket))
		err := b.Put([]byte(uuid), []byte(extension))
		return err
	})
}

func (d *database) FindExtension(uuid string) string {
	var extension string
	d.bolt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(ExtensionBucket))
		e := b.Get([]byte(uuid))
		extension = string(e)
		return nil
	})

	return extension
}
