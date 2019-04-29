package Indexer

import (
	"fmt"
	"os"
	"strconv"

	"github.com/dgraph-io/badger"
)

// URL -> Page ID Indexer and Word -> Page ID Indexer
type PageRankIndexer struct {
	db           *badger.DB
	databasePath string
}

// After initializing the PageRankIndexer, we need to call defer PageRankIndexer.Release()
func (pageRankIndexer *PageRankIndexer) Initialize(path string) error {
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
	pageRankIndexer.db = db
	pageRankIndexer.databasePath = path
	return err
}

func (pageRankIndexer *PageRankIndexer) Release() error {
	return pageRankIndexer.db.Close()
}

func (pageRankIndexer *PageRankIndexer) Backup() error {
	fmt.Println("Doing Database Backup")
	f, err := os.Create(pageRankIndexer.databasePath)
	if err != nil {
		return err
	}
	defer f.Close()
	pageRankIndexer.db.Backup(f, 0)
	return err
}

func (pageRankIndexer *PageRankIndexer) AddKeyToIndex(key uint64, value float64) error {
	var err error
	err = pageRankIndexer.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(uint64ToByte(key))
		if err == badger.ErrKeyNotFound {
			// Get new value for index
			v := fmt.Sprintf("%.6f", value)
			err = txn.Set(uint64ToByte(key), []byte(v))
			return err
		}
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error in adding Key to Index: %s", err)
	}
	return err
}

func (pageRankIndexer *PageRankIndexer) GetValueFromKey(key uint64) (float64, error) {
	var result float64
	var floatErr error
	err := pageRankIndexer.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(uint64ToByte(key))
		if err == nil {
			itemErr := item.Value(func(val []byte) error {
				result, floatErr = strconv.ParseFloat(string(val), 64)
				if floatErr != nil {
					return floatErr
				}
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

func (pageRankIndexer *PageRankIndexer) Iterate() {
	fmt.Println("Iterating over Mapping Index")
	_ = pageRankIndexer.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				floatValue, floatErr := strconv.ParseFloat(string(v), 64)
				if floatErr != nil {
					return floatErr
				}
				fmt.Printf("key=%d, value=%f\n", byteToUint64(k), floatValue)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (pageRankIndexer *PageRankIndexer) DeleteKeyValuePair(key uint64) error {
	err := pageRankIndexer.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(uint64ToByte(key))
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error when deleting value from key: %s", err)
	}
	return err
}
