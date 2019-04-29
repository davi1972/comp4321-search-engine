package Indexer

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
)

type VSMIndexer struct {
	db           *badger.DB
	databasePath string
}

// After initializing the VSMIndexer, we need to call defer VSMIndexer.Release()
func (VSMIndexer *VSMIndexer) Initialize(path string) error {
	if err := os.MkdirAll(path, 0774); err != nil {
		return err
	}
	opts := badger.DefaultOptions
	opts.Dir = path
	opts.ValueDir = path
	db, err := badger.Open(opts)
	if err != nil {
		return fmt.Errorf("Error while initializing: %s", err)
	}
	VSMIndexer.db = db
	VSMIndexer.databasePath = path
	return err
}

func (VSMIndexer *VSMIndexer) Release() error {
	return VSMIndexer.db.Close()
}

func (VSMIndexer *VSMIndexer) Backup() error {
	fmt.Println("Doing Database Backup")
	f, err := os.Create(VSMIndexer.databasePath)
	if err != nil {
		return err
	}
	defer f.Close()
	VSMIndexer.db.Backup(f, 0)
	return err
}

func (VSMIndexer *VSMIndexer) AddKeyToIndex(key uint64, value uint64) error {
	var err error
	err = VSMIndexer.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(uint64ToByte(key))
		if err == badger.ErrKeyNotFound {
			err = txn.Set(uint64ToByte(key), uint64ToByte(value))
			return err
		}
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error in adding Key to Index: %s", err)
	}
	return err
}

func (VSMIndexer *VSMIndexer) GetValueFromKey(key uint64) (uint64, error) {
	var result uint64
	err := VSMIndexer.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(uint64ToByte(key))
		if err == nil {
			itemErr := item.Value(func(val []byte) error {
				result = byteToUint64(val)
				return nil
			})
			if itemErr != nil {
				return itemErr
			}
		}
		return err
	})

	if err != nil {
		err = fmt.Errorf("Error in getting Value from Key: %s", err)
	}
	return result, err
}

func (VSMIndexer *VSMIndexer) Iterate() {
	fmt.Println("Iterating over Mapping Index")
	_ = VSMIndexer.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				fmt.Printf("key=%d, value=%d\n", byteToUint64(k), byteToUint64(v))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (VSMIndexer *VSMIndexer) DeleteKeyValuePair(key uint64) error {
	err := VSMIndexer.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(uint64ToByte(key))
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error when deleting value from key: %s", err)
	}
	return err
}
