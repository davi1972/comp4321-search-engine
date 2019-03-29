package Indexer

import (
	"github.com/dgraph-io/badger"
	"os"
	"fmt"
	"strings"
	"strconv"
)

type PagePropetiesIndexer struct {
	db *badger.DB
	databasePath string
}

type Page struct {
	id uint64
	title string
	url string
}

func CreatePage(id uint64, title string, url string) Page {
	page := Page{}
	page.id = id
	page.title = title
	page.url = url
	return page
}

func pageToString(page *Page) string {
	return string(uint64ToByte(page.id)) + "/page/" + page.title + "/page/" + page.url
}

func stringToPage(str string) Page {
	splitString := strings.Split(str, "/page/")
	idString, _ := strconv.ParseUint(splitString[0], 10, 64)
	return Page{idString, splitString[1], splitString[2]}
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

func (PagePropetiesIndexer *PagePropetiesIndexer) AddKeyToPageProperties(pageID uint64, page Page) error {
	pageString := uint64ToByte(pageID)
	err := PagePropetiesIndexer.db.Update(func(txn *badger.Txn) error {
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

func (PagePropetiesIndexer *PagePropetiesIndexer) GetPagePropertiesFromKey(pageID uint64) (Page, error) {
	pageString := uint64ToByte(pageID)
	var resultPage Page
	err := PagePropetiesIndexer.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(pageString)
		if err != nil {
			return err
		} 
		itemErr := item.Value(func(val []byte) error {
			fmt.Println("Get value: %s", string(val))
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


func (PagePropetiesIndexer *PagePropetiesIndexer) DeletePagePropertiesFromKey(pageID uint64) error {
	pageString := uint64ToByte(pageID)
	err := PagePropetiesIndexer.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(pageString))
		return err
	})
	if err != nil {
		err = fmt.Errorf("Error when deleting page properties from key: %s", err)
	}
	return err
}