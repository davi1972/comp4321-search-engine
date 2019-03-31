package Indexer

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/dgraph-io/badger"
)

type InvertedFileIndexer struct {
	db           *badger.DB
	databasePath string
}

type InvertedFile struct {
	pageID        uint64
	wordPositions []uint64
}

func CreateInvertedFile(pageID uint64) *InvertedFile {
	return &InvertedFile{pageID, []uint64{}}
}

func (invertedFile *InvertedFile) AddWordPositions(pos uint64) {
	invertedFile.wordPositions = append(invertedFile.wordPositions, pos)
}

func (invertedFileIndexer *InvertedFileIndexer) Initialize(path string) error {
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

	invertedFileIndexer.db = db
	invertedFileIndexer.databasePath = path
	return nil
}

func (invertedFileIndexer *InvertedFileIndexer) Release() error {
	return invertedFileIndexer.db.Close()
}

func (invertedFileIndexer *InvertedFileIndexer) Iterate() error {
	fmt.Println("iterating over InvertedFile")
	err := invertedFileIndexer.db.View(func(txn *badger.Txn) error {
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

func stringToInvertedFile(str string) InvertedFile {
	pageSliceString := strings.Split(str, " ")
	pageSliceUint64 := make([]uint64, len(pageSliceString))
	for i, value := range pageSliceString {
		pageSliceUint64[i], _ = strconv.ParseUint(value, 10, 64)
	}
	return InvertedFile{pageSliceUint64[0], pageSliceUint64[1:]}
}

func (InvertedFileIndexer *InvertedFileIndexer) Backup() error {
	fmt.Println("Doing Database Backup")
	f, err := os.Create(InvertedFileIndexer.databasePath)
	if err != nil {
		return err
	}
	defer f.Close()
	InvertedFileIndexer.db.Backup(f, 0)
	return err
}

func invertedFileToString(i InvertedFile) string {
	result := strconv.FormatUint(i.pageID, 10)
	wordPositionsString := make([]string, len(i.wordPositions))
	for i, v := range i.wordPositions {
		wordPositionsString[i] = strconv.FormatUint(v, 10)
	}
	return result + " " + strings.Join(wordPositionsString, " ")
}

func (invertedFileIndexer *InvertedFileIndexer) AddKeyToIndexOrUpdate(wordID uint64, invertedFile InvertedFile) error {
	keyString := uint64ToByte(wordID)

	// Construct a string to to add to inverted file
	valueString := invertedFileToString(invertedFile)
	err := invertedFileIndexer.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get(keyString)
		// If key already exists, have to append/insert
		if err == nil {
			itemErr := item.Value(func(val []byte) error {
				// First convert to inverted file slice
				invertedFileListString := strings.Split(string(val), ",")
				invertedFileList := make([]InvertedFile, len(invertedFileListString))
				for i, v := range invertedFileListString {
					invertedFileList[i] = stringToInvertedFile(v)
				}

				// Special case if the inverted file is the largest
				if invertedFileList[len(invertedFileList)-1].pageID < invertedFile.pageID {
					invertedFileList = append(invertedFileList, invertedFile)
				} else {
					// Insert to sorted Inverted File String
					for i, v := range invertedFileList {
						if v.pageID < invertedFile.pageID {
							invertedFileList = append(invertedFileList, InvertedFile{})
							copy(invertedFileList[i+1:], invertedFileList[i:])
							invertedFileList[i] = invertedFile
							break
						}

					}
				}
				valueString = invertedFileToString(invertedFileList[0])
				for _, v := range invertedFileList[1:] {
					valueString = valueString + "," + invertedFileToString(v)
				}
				return nil
			})
			if itemErr != nil {
				return itemErr
			}

			// Delete the old one
			err = txn.Delete(keyString)
		}
		err = txn.Set(keyString, []byte(valueString))
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error in adding Key to Index or Update: %s", err)
	}
	return err
}

func (invertedFileIndexer *InvertedFileIndexer) GetInvertedFileFromKey(wordID uint64) ([]InvertedFile, error) {
	keyString := uint64ToByte(wordID)
	result := make([]InvertedFile, 0)
	err := invertedFileIndexer.db.View(func(txn *badger.Txn) error {
		item, getErr := txn.Get(keyString)
		if getErr == nil {
			_ = item.Value(func(val []byte) error {
				resultString := string(val)
				resultList := strings.Split(resultString, ",")
				for _, v := range resultList {
					result = append(result, stringToInvertedFile(v))
				}
				return nil
			})
		}
		return getErr
	})
	if err != nil {
		err = fmt.Errorf("Error when getting value transaction: %s", err)
	}
	return result, err
}

func (invertedFileIndexer *InvertedFileIndexer) DeleteAllInvertedFileFromKey(wordID uint64) error {
	keyString := uint64ToByte(wordID)
	err := invertedFileIndexer.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(keyString))
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error when deleting key value pair: %s", err)
	}
	return err
}

func (invertedFileIndexer *InvertedFileIndexer) DeleteInvertedFileFromWordListAndPage(wordIDList []uint64, pageID uint64) error {
	// pageString := uint64ToByte(pageID)
	var err error
	fmt.Println("in delete")
	for _, word := range wordIDList {
		fmt.Printf("Using Word: %d \n", word)
		err = invertedFileIndexer.db.Update(func(txn *badger.Txn) error {
			item, err := txn.Get(uint64ToByte(word))
			if err == badger.ErrKeyNotFound {
				return err
			} else if err == nil {
				err = item.Value(func(val []byte) error {
					resultString := string(val)
					resultList := strings.Split(resultString, ",")
					invertedFileListString := make([][]string, 0)
					for _, v := range resultList {
						invertedFileListString = append(invertedFileListString, strings.Split(v, " "))
					}
					fmt.Println("Going to Delete:: ")
					fmt.Println(resultString)
					fmt.Println(invertedFileListString)
					return nil
				})
			}

			return err
		})
	}
	return err
}
