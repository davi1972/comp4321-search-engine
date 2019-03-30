package Indexer

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dgraph-io/badger"
)

type DocumenWordForwardIndexer struct {
	db           *badger.DB
	databasePath string
}

type WordFrequency struct {
	wordID    uint64
	frequency uint64
}

func CreateWordFrequency(id uint64, f uint64) WordFrequency {
	return WordFrequency{id, f}
}

func wordFrequencyToString(word *WordFrequency) string {
	return strconv.Itoa(int(word.wordID)) + " " + strconv.Itoa(int(word.frequency))
}

func stringToWordFrequency(str string) WordFrequency {
	splitString := strings.Split(str, "/page/")
	id, _ := strconv.Atoi(splitString[0])
	freq, _ := strconv.Atoi(splitString[1])
	return WordFrequency{uint64(id), uint64(freq)}
}

func (documenWordForwardIndexer *DocumenWordForwardIndexer) Initialize(path string) error {
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
	documenWordForwardIndexer.db = db
	documenWordForwardIndexer.databasePath = path
	return err
}

func (documenWordForwardIndexer *DocumenWordForwardIndexer) Release() error {
	return documenWordForwardIndexer.db.Close()
}

func (documenWordForwardIndexer *DocumenWordForwardIndexer) Backup() error {
	fmt.Println("Doing Database Backup")
	f, err := os.Create(documenWordForwardIndexer.databasePath)
	if err != nil {
		return err
	}
	defer f.Close()
	documenWordForwardIndexer.db.Backup(f, 0)
	return err
}

// func (documenWordForwardIndexer *DocumenWordForwardIndexer) AddWordFrequencyListToKey(documentId uint64, wordFrequencyList []WordFrequency) error {
// 	var valueString string
// 	if len(wordFrequencyList) > 0 {
// 		valueString = wordFrequencyToString(&wordFrequencyList[0])
// 		for _, word := range wordFrequencyList {

// 		}
// 	}
// 	err := documenWordForwardIndexer.db.Update(func(txn *badger.Txn) error {
// 		err := txn.Set(uint64ToByte(documentId), []byte(valueString))
// 		return err
// 	})
// 	if err != nil {
// 		err = fmt.Errorf("Error in adding Key to Index: %s", err)
// 	}
// 	return err
// }

func (documenWordForwardIndexer *DocumenWordForwardIndexer) GetIdListFromKey(documentId uint64) ([]uint64, error) {
	result := make([]uint64, 0)
	err := documenWordForwardIndexer.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(uint64ToByte(documentId))
		if err == nil {
			itemErr := item.Value(func(val []byte) error {
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

func (documenWordForwardIndexer *DocumenWordForwardIndexer) DeleteKeyValuePair(documentId uint64) error {
	err := documenWordForwardIndexer.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(uint64ToByte(documentId))
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error when deleting value from key: %s", err)
	}
	return err
}

func (documenWordForwardIndexer *DocumenWordForwardIndexer) Iterate() error {
	fmt.Println("iterating over Forward Index")
	err := documenWordForwardIndexer.db.View(func(txn *badger.Txn) error {
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
