package Indexer

import (
	"github.com/dgraph-io/badger"
	"fmt"
)

// URL -> Page ID Indexer and Word -> Page ID Indexer
type MappingIndexer struct {
	db *badger.DB
	sequence *badger.Sequence
}

// After initializing the mappingIndexer, we need to call defer mappingIndexer.Release()
func (mappingIndexer *MappingIndexer) Initialize() error {
	opts := badger.DefaultOptions
	opts.Dir = "../tmp/mappingdb"
	opts.ValueDir = "../tmp/mappingdb"
	db, err := badger.Open(opts)
	if err != nil {
		return fmt.Errorf("Error while initializing: %s", err)
	}
	mappingIndexer.db = db 
	mappingIndexer.sequence, _ = db.GetSequence([]byte(""), 1000)
	return nil
}

func (mappingIndexer *MappingIndexer) Release() error {
	mappingIndexer.sequence.Release()
	return mappingIndexer.db.Close()
}

func (mappingIndexer *MappingIndexer) AddKeyToIndex(key string) error {
	err := mappingIndexer.db.Update(func(txn *badger.Txn) error {
		// Get new value for index
		id, err := mappingIndexer.sequence.Next()

		err = txn.Set([]byte(key), []byte(uint64ToByte(id)))
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error in adding Key to Index: %s", err)
	}
	return err
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