package Indexer

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger"
)

//  Word ID -> word Indexer
type ReverseMappingIndexer struct {
	db           *badger.DB
	databasePath string
}

// After initializing the ReverseMappingIndexer, we need to call defer ReverseMappingIndexer.Release()
func (reverseMappingIndexer *ReverseMappingIndexer) Initialize(path string) error {
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
	reverseMappingIndexer.db = db
	reverseMappingIndexer.databasePath = path
	return err
}

func (reverseMappingIndexer *ReverseMappingIndexer) Release() error {
	return reverseMappingIndexer.db.Close()
}

func (reverseMappingIndexer *ReverseMappingIndexer) Backup() error {
	fmt.Println("Doing Database Backup")
	f, err := os.Create(reverseMappingIndexer.databasePath)
	if err != nil {
		return err
	}
	defer f.Close()
	reverseMappingIndexer.db.Backup(f, 0)
	return err
}

func (reverseMappingIndexer *ReverseMappingIndexer) AddKeyToIndex(wordID uint64, word string) error {
	err := reverseMappingIndexer.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(uint64ToByte(wordID))
		if err == badger.ErrKeyNotFound {
			// Get new value for index
			err = txn.Set(uint64ToByte(wordID), []byte(word))
			return err
		}
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error in adding Key to Index: %s", err)
	}
	return err
}

func (ReverseMappingIndexer *ReverseMappingIndexer) GetValueFromKey(wordID uint64) (string, error) {
	var result string
	err := ReverseMappingIndexer.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(uint64ToByte(wordID))
		if err == nil {
			itemErr := item.Value(func(val []byte) error {
				result = string(val)
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

func (ReverseMappingIndexer *ReverseMappingIndexer) Iterate() {
	fmt.Println("Iterating over reverse Mapping Index")
	_ = ReverseMappingIndexer.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				fmt.Printf("key=%d, value=%s\n", byteToUint64(k), string(v))
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (ReverseMappingIndexer *ReverseMappingIndexer) DeleteKeyValuePair(wordID uint64) error {
	err := ReverseMappingIndexer.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(uint64ToByte(wordID))
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error when deleting value from key: %s", err)
	}
	return err
}
