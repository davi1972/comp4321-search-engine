package Indexer

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dgraph-io/badger"
)

type DocumentWordForwardIndexer struct {
	db           *badger.DB
	databasePath string
}

type WordFrequency struct {
	wordID    uint64
	frequency uint64
}

func (w *WordFrequency) GetID() uint64 {
	return w.wordID
}

func (w *WordFrequency) GetFrequency() uint64 {
	return w.frequency
}

func CreateWordFrequency(id uint64, f uint64) WordFrequency {
	return WordFrequency{id, f}
}

func wordFrequencyToString(word *WordFrequency) string {
	return strconv.Itoa(int(word.wordID)) + " " + strconv.Itoa(int(word.frequency))
}

func stringToWordFrequency(str string) WordFrequency {
	splitString := strings.Split(str, " ")
	id, _ := strconv.Atoi(splitString[0])
	freq, _ := strconv.Atoi(splitString[1])
	return WordFrequency{uint64(id), uint64(freq)}
}

func (w *WordFrequency) GetWordID() uint64 {
	return w.wordID
}

func (documentWordForwardIndexer *DocumentWordForwardIndexer) Initialize(path string) error {
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
	documentWordForwardIndexer.db = db
	documentWordForwardIndexer.databasePath = path
	return err
}

func (documentWordForwardIndexer *DocumentWordForwardIndexer) Release() error {
	return documentWordForwardIndexer.db.Close()
}

func (documentWordForwardIndexer *DocumentWordForwardIndexer) Backup() error {
	fmt.Println("Doing Database Backup")
	f, err := os.Create(documentWordForwardIndexer.databasePath)
	if err != nil {
		return err
	}
	defer f.Close()
	documentWordForwardIndexer.db.Backup(f, 0)
	return err
}

func (documentWordForwardIndexer *DocumentWordForwardIndexer) AddWordFrequencyListToKey(documentId uint64, wordFrequencyList []WordFrequency) error {
	var valueString string
	if len(wordFrequencyList) > 0 {
		valueString = wordFrequencyToString(&wordFrequencyList[0])
		for _, word := range wordFrequencyList[1:] {
			valueString = valueString + "," + wordFrequencyToString(&word)
		}
	}
	err := documentWordForwardIndexer.db.Update(func(txn *badger.Txn) error {
		err := txn.Set(uint64ToByte(documentId), []byte(valueString))
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error in adding Key to Index: %s", err)
	}
	return err
}

func (documentWordForwardIndexer *DocumentWordForwardIndexer) GetWordFrequencyListFromKey(documentId uint64) ([]WordFrequency, error) {
	result := make([]WordFrequency, 0)
	err := documentWordForwardIndexer.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(uint64ToByte(documentId))
		if err == nil {
			itemErr := item.Value(func(val []byte) error {
				if string(val) != "" {

					resultList := strings.Split(string(val), ",")
					for _, v := range resultList {
						result = append(result, stringToWordFrequency(v))
					}
					return nil
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

func (documentWordForwardIndexer *DocumentWordForwardIndexer) DeleteKeyValuePair(documentId uint64) error {
	err := documentWordForwardIndexer.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete(uint64ToByte(documentId))
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error when deleting value from key: %s", err)
	}
	return err
}

func (documentWordForwardIndexer *DocumentWordForwardIndexer) Iterate() error {
	fmt.Println("iterating over Document Word Forward Index")
	err := documentWordForwardIndexer.db.View(func(txn *badger.Txn) error {
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

// find N = num of docs
func (documentWordForwardIndexer *DocumentWordForwardIndexer) GetSize() uint64 {
	//fmt.Println("Iterating over Document Word Forward Index to count size")
	i := 0
	_ = documentWordForwardIndexer.db.View(func(txn *badger.Txn) error {
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

func (documentWordForwardIndexer *DocumentWordForwardIndexer) GetDocIDList() ([]uint64, error) {
	//fmt.Println("Iterating over Document Word Forward Index for Doc IDs")
	result := make([]uint64, 0)
	err := documentWordForwardIndexer.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				result = append(result, byteToUint64(k))
				return nil
			})
			if err != nil {
				return nil
			}
		}
		return nil
	})
	// fmt.Printf("Size of doc ID List: %d\n", len(result))
	// fmt.Printf("Values in result: %v\n", result)
	return result, err
}
