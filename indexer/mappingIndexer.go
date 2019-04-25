package Indexer

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
)

// URL -> Page ID Indexer and Word -> Page ID Indexer
type MappingIndexer struct {
	db           *badger.DB
	sequence     *badger.Sequence
	databasePath string
}

// After initializing the mappingIndexer, we need to call defer mappingIndexer.Release()
func (mappingIndexer *MappingIndexer) Initialize(path string) error {
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
	mappingIndexer.db = db
	mappingIndexer.sequence, _ = db.GetSequence([]byte("doc-"), 1000)
	mappingIndexer.databasePath = path
	return err
}

func (mappingIndexer *MappingIndexer) Release() error {
	mappingIndexer.sequence.Release()
	return mappingIndexer.db.Close()
}

func (mappingIndexer *MappingIndexer) Backup() error {
	fmt.Println("Doing Database Backup")
	f, err := os.Create(mappingIndexer.databasePath)
	if err != nil {
		return err
	}
	defer f.Close()
	mappingIndexer.db.Backup(f, 0)
	return err
}

func (mappingIndexer *MappingIndexer) AddKeyToIndex(key string) (uint64, error) {
	var id uint64
	var err error
	err = mappingIndexer.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err == badger.ErrKeyNotFound {
			// Get new value for index
			id, err = mappingIndexer.sequence.Next()
			err = txn.Set([]byte(key), []byte(uint64ToByte(id)))
			return err
		}
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error in adding Key to Index: %s", err)
	}
	return id, err
}

func (mappingIndexer *MappingIndexer) GetValueFromKey(key string) (uint64, error) {
	var result uint64
	err := mappingIndexer.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
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

func (mappingIndexer *MappingIndexer) Iterate() {
	fmt.Println("Iterating over Mapping Index")
	_ = mappingIndexer.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				fmt.Printf("key=%s, value=%d\n", k, byteToUint64(v))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (mappingIndexer *MappingIndexer) DeleteKeyValuePair(key string) error {
	err := mappingIndexer.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error when deleting value from key: %s", err)
	}
	return err
}

// find N = num of docs
func (mappingIndexer *MappingIndexer) GetSize() uint64 {
	fmt.Println("Iterating over Mapping Index to count size")
	i := 0
	_ = mappingIndexer.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			i++
		}
		return nil
	})
	return uint64(i)
}
