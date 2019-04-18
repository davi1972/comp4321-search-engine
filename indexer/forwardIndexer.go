package Indexer

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dgraph-io/badger"
)

type ForwardIndexer struct {
	db           *badger.DB
	databasePath string
}

func (forwardIndexer *ForwardIndexer) Initialize(path string) error {
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
	forwardIndexer.db = db
	forwardIndexer.databasePath = path
	return err
}

func (forwardIndexer *ForwardIndexer) Release() error {
	return forwardIndexer.db.Close()
}

func (forwardIndexer *ForwardIndexer) Backup() error {
	fmt.Println("Doing Database Backup")
	f, err := os.Create(forwardIndexer.databasePath)
	if err != nil {
		return err
	}
	defer f.Close()
	forwardIndexer.db.Backup(f, 0)
	return err
}

func (forwardIndexer *ForwardIndexer) AddIdListToKey(documentId uint64, idList []uint64) error {
	var valueString string
	if len(idList) > 0 {
		valueString = strconv.FormatUint(idList[0], 10)
		if len(idList) > 1 {
			for _, v := range idList[1:] {
				valueString = valueString + " " + strconv.FormatUint(v, 10)
			}
		}
	}
	err := forwardIndexer.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(uint64ToByte(documentId), []byte(valueString))
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error in adding Key to Index: %s", err)
	}
	return err
}

func (forwardIndexer *ForwardIndexer) GetIdListFromKey(documentId uint64) ([]uint64, error) {
	result := make([]uint64, 0)
	err := forwardIndexer.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(uint64ToByte(documentId))
		if err == nil {
			itemErr := item.Value(func(val []byte) error {
				if string(val) == "" {
					return nil
				}
				resultList := strings.Split(string(val), " ")
				for _, v := range resultList {
					val, _ := strconv.Atoi(v)
					result = append(result, uint64(val))
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

func (forwardIndexer *ForwardIndexer) DeleteKeyValuePair(documentId uint64) error {
	err := forwardIndexer.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(uint64ToByte(documentId))
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error when deleting value from key: %s", err)
	}
	return err
}

func (forwardIndexer *ForwardIndexer) Iterate() error {
	fmt.Println("Iterating over Forward Index")
	err := forwardIndexer.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				fmt.Printf("key=%d, value=%s\n", byteToUint64(k), v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
