package Indexer

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgraph-io/badger"
)

type PagePropetiesIndexer struct {
	db           *badger.DB
	databasePath string
}

type Page struct {
	id           uint64
	title        string
	url          string
	size         int
	dateModified time.Time
}

func (page *Page) GetId() uint64 {
	return page.id
}

func (page *Page) GetTitle() string {
	return page.title
}

func (page *Page) GetUrl() string {
	return page.url
}

func (page *Page) GetSize() int {
	return page.size
}

func (page *Page) GetDateString() string {
	return page.dateModified.Format(time.RFC1123)
}

func (page *Page) GetDate() time.Time {
	return page.dateModified
}

func CreatePage(id uint64, title string, url string, size int, date time.Time) Page {
	page := Page{}
	page.id = id
	page.title = title
	page.url = url
	page.size = size
	page.dateModified = date
	return page
}

func pageToString(page *Page) string {
	return strconv.Itoa(int(page.id)) + "/page/" + page.title + "/page/" + page.url + "/page/" + strconv.FormatInt(int64(page.size), 10) + "/page/" + page.dateModified.Format(time.RFC3339)
}

func stringToPage(str string) Page {
	splitString := strings.Split(str, "/page/")
	idString, _ := strconv.ParseUint(splitString[0], 10, 64)
	size, _ := strconv.Atoi(splitString[3])
	time, _ := time.Parse(time.RFC3339, splitString[4])
	return Page{idString, splitString[1], splitString[2], size, time}
}

// After initializing the PagePropetiesIndexer, we need to call defer PagePropetiesIndexer.Release()
func (pagePropetiesIndexer *PagePropetiesIndexer) Initialize(path string) error {
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
	pagePropetiesIndexer.db = db
	pagePropetiesIndexer.databasePath = path
	return err
}

func (pagePropetiesIndexer *PagePropetiesIndexer) All() ([]Page, error) {

	pages := []Page{}

	err := pagePropetiesIndexer.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			var p Page
			item := it.Item()
			// k := item.Key()
			err := item.Value(func(v []byte) error {
				p = stringToPage(string(v))
				return nil
			})
			if err != nil {
				return err
			}

			pages = append(pages, p)

		}
		return nil
	})

	return pages, err
}

func (pagePropetiesIndexer *PagePropetiesIndexer) Iterate() error {
	fmt.Println("iterating over Page Properties")
	err := pagePropetiesIndexer.db.View(func(txn *badger.Txn) error {
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

func (pagePropetiesIndexer *PagePropetiesIndexer) Release() error {
	return pagePropetiesIndexer.db.Close()
}

func (pagePropetiesIndexer *PagePropetiesIndexer) Backup() error {
	fmt.Println("Doing Database Backup")
	f, err := os.Create(pagePropetiesIndexer.databasePath)
	if err != nil {
		return err
	}
	defer f.Close()
	pagePropetiesIndexer.db.Backup(f, 0)
	return err
}

func (pagePropetiesIndexer *PagePropetiesIndexer) AddKeyToPageProperties(pageID uint64, page Page) error {
	pageString := uint64ToByte(pageID)
	err := pagePropetiesIndexer.db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get(pageString)

		// If key already exists, have to delete and add new one
		if err == nil {
			txn.Delete([]byte(pageString))
		}

		pagePropetiesString := pageToString(&page)

		err = txn.Set(uint64ToByte(pageID), []byte(pagePropetiesString))
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error in adding Key to Index: %s", err)
	}
	return err
}

func (pagePropetiesIndexer *PagePropetiesIndexer) GetPagePropertiesFromKey(pageID uint64) (Page, error) {
	pageString := uint64ToByte(pageID)
	var resultPage Page
	err := pagePropetiesIndexer.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(pageString)
		if err != nil {
			return err
		}
		itemErr := item.Value(func(val []byte) error {
			resultPage = stringToPage(string(val))
			return nil
		})
		if itemErr != nil {
			return itemErr
		}
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error when getting page properties from key: %s", err)
	}
	return resultPage, err
}

func (pagePropetiesIndexer *PagePropetiesIndexer) DeletePagePropertiesFromKey(pageID uint64) error {
	pageString := uint64ToByte(pageID)
	err := pagePropetiesIndexer.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(pageString))
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error when deleting page properties from key: %s", err)
	}
	return err
}
