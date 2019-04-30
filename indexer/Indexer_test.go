package Indexer

import (
	"os"
	"testing"
	"time"
)

func TestInitializeDocumentWordForwardIndexer(t *testing.T) {

	wd, _ := os.Getwd()
	testDB := &DocumentWordForwardIndexer{}
	err := testDB.Initialize(wd + "/dbTest/documentWordForwardIndex")
	defer testDB.Release()
	if err != nil {
		t.Fail()
	}
}

func TestIterateDocumentWordForwardIndexer(t *testing.T) {

	wd, _ := os.Getwd()
	testDB := &DocumentWordForwardIndexer{}
	err := testDB.Initialize(wd + "/dbTest/documentWordForwardIndex")
	defer testDB.Release()
	if err != nil {
		t.Fail()
	}
	testDB.Iterate()
}

func TestAddGetDatabaseDocumentWordForwardIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &DocumentWordForwardIndexer{}
	err := testDB.Initialize(wd + "/dbTest/documentWordForwardIndex")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}

	wordF := make([]WordFrequency, 1)
	wordF[0].Frequency = 10
	wordF[0].WordID = 0

	addErr := testDB.AddWordFrequencyListToKey(0, wordF)
	if addErr != nil {
		t.FailNow()
	}

	wordFResult, resultErr := testDB.GetWordFrequencyListFromKey(0)
	if resultErr != nil {
		t.FailNow()
	}

	if wordFResult[0].Frequency != wordF[0].Frequency || wordFResult[0].WordID != wordF[0].WordID {
		t.Fail()
	}
}

func TestGetSizeDocumentWordForwardIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &DocumentWordForwardIndexer{}
	err := testDB.Initialize(wd + "/dbTest/documentWordForwardIndex")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}

	wordF := make([]WordFrequency, 1)
	wordF[0].Frequency = 10
	wordF[0].WordID = 0

	addErr := testDB.AddWordFrequencyListToKey(0, wordF)
	if addErr != nil {
		t.FailNow()
	}
	size := testDB.GetSize()

	if size != 1 {
		t.Fail()
	}
}

func TestGetDocIDListDocumentWordForwardIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &DocumentWordForwardIndexer{}
	err := testDB.Initialize(wd + "/dbTest/documentWordForwardIndex")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}

	wordF := make([]WordFrequency, 1)
	wordF[0].Frequency = 10
	wordF[0].WordID = 0

	addErr := testDB.AddWordFrequencyListToKey(0, wordF)
	if addErr != nil {
		t.FailNow()
	}
	docIDList, docIDErr := testDB.GetDocIDList()
	if docIDErr != nil {
		t.Fail()
	}
	if len(docIDList) != 1 {
		t.Fail()
	}
}

func TestDeleteDatabaseDocumentWordForwardIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &DocumentWordForwardIndexer{}
	err := testDB.Initialize(wd + "/dbTest/documentWordForwardIndex")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}

	wordF := make([]WordFrequency, 1)
	wordF[0].Frequency = 10
	wordF[0].WordID = 0

	addErr := testDB.AddWordFrequencyListToKey(0, wordF)
	if addErr != nil {
		t.FailNow()
	}

	deleteErr := testDB.DeleteKeyValuePair(0)
	if deleteErr != nil {
		t.FailNow()
	}

	_, resultErr := testDB.GetWordFrequencyListFromKey(0)
	if resultErr == nil {
		t.FailNow()
	}

}

func TestInitializeForwardIndexer(t *testing.T) {

	wd, _ := os.Getwd()
	testDB := &ForwardIndexer{}
	err := testDB.Initialize(wd + "/dbTest/ForwardIndex")
	defer testDB.Release()
	if err != nil {
		t.Fail()
	}
}

func TestIterateForwardIndexer(t *testing.T) {

	wd, _ := os.Getwd()
	testDB := &ForwardIndexer{}
	err := testDB.Initialize(wd + "/dbTest/ForwardIndexer")
	defer testDB.Release()
	if err != nil {
		t.Fail()
	}
	testDB.Iterate()
}

func TestAddGetDatabaseForwardIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &ForwardIndexer{}
	err := testDB.Initialize(wd + "/dbTest/ForwardIndex")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}

	idList := make([]uint64, 1)
	idList[0] = 0

	addErr := testDB.AddIdListToKey(0, idList)
	if addErr != nil {
		t.FailNow()
	}

	idListResult, resultErr := testDB.GetIdListFromKey(0)
	if resultErr != nil {
		t.FailNow()
	}

	if idListResult[0] != 0 {
		t.Fail()
	}
}

func TestDeleteDatabaseForwardIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &ForwardIndexer{}
	err := testDB.Initialize(wd + "/dbTest/ForwardIndex")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}

	idList := make([]uint64, 1)
	idList[0] = 0

	addErr := testDB.AddIdListToKey(0, idList)
	if addErr != nil {
		t.FailNow()
	}

	deleteErr := testDB.DeleteKeyValuePair(0)
	if deleteErr != nil {
		t.FailNow()
	}

	_, resultErr := testDB.GetIdListFromKey(0)
	if resultErr == nil {
		t.FailNow()
	}

}

func TestInitializeMappingIndexer(t *testing.T) {

	wd, _ := os.Getwd()
	testDB := &MappingIndexer{}
	err := testDB.Initialize(wd + "/dbTest/MappingIndex")
	defer testDB.Release()
	if err != nil {
		t.Fail()
	}
}

func TestIterateMappingIndexer(t *testing.T) {

	wd, _ := os.Getwd()
	testDB := &MappingIndexer{}
	err := testDB.Initialize(wd + "/dbTest/MappingIndexer")
	defer testDB.Release()
	if err != nil {
		t.Fail()
	}
	testDB.Iterate()
}

func TestAddGetDatabaseMappingIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &MappingIndexer{}
	err := testDB.Initialize(wd + "/dbTest/MappingIndex")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}

	id, addErr := testDB.AddKeyToIndex("test")
	if addErr != nil {
		t.FailNow()
	}

	idResult, resultErr := testDB.GetValueFromKey("test")
	if resultErr != nil {
		t.FailNow()
	}

	if idResult != id {
		t.Fail()
	}
}

func TestAllDatabaseMappingIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &MappingIndexer{}
	err := testDB.Initialize(wd + "/dbTest/MappingIndex")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}

	_, addErr := testDB.AddKeyToIndex("test")
	if addErr != nil {
		t.FailNow()
	}

	_, keysErr := testDB.All()
	if keysErr != nil {
		t.FailNow()
	}

}

func TestDeleteDatabaseMappingIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &MappingIndexer{}
	err := testDB.Initialize(wd + "/dbTest/MappingIndex")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}

	idList := make([]uint64, 1)
	idList[0] = 0

	_, addErr := testDB.AddKeyToIndex("test")
	if addErr != nil {
		t.FailNow()
	}

	deleteErr := testDB.DeleteKeyValuePair("test")
	if deleteErr != nil {
		t.FailNow()
	}

	_, resultErr := testDB.GetValueFromKey("test")
	if resultErr == nil {
		t.FailNow()
	}

}

func TestInitializeInvertedFileIndexer(t *testing.T) {

	wd, _ := os.Getwd()
	testDB := &InvertedFileIndexer{}
	err := testDB.Initialize(wd + "/dbTest/InvertedFileIndexer")
	defer testDB.Release()
	if err != nil {
		t.Fail()
	}
}

func TestIterateInvertedFileIndexer(t *testing.T) {

	wd, _ := os.Getwd()
	testDB := &InvertedFileIndexer{}
	err := testDB.Initialize(wd + "/dbTest/InvertedFileIndexer")
	defer testDB.Release()
	if err != nil {
		t.Fail()
	}
	testDB.Iterate()
}

func TestDeleteWholeKeyInvertedFileIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &InvertedFileIndexer{}
	err := testDB.Initialize(wd + "/dbTest/InvertedFileIndexer")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}

	invFile := InvertedFile{pageID: 0}
	invFile.AddWordPositions(0)
	invFile.AddWordPositions(1)

	addErr := testDB.AddKeyToIndexOrUpdate(0, invFile)
	if addErr != nil {
		t.FailNow()
	}

	invFile2 := InvertedFile{pageID: 1}
	invFile2.AddWordPositions(2)
	invFile2.AddWordPositions(3)

	addErr = testDB.AddKeyToIndexOrUpdate(0, invFile2)
	if addErr != nil {
		t.FailNow()
	}

	testDB.DeleteAllInvertedFileFromKey(0)

	invFileResult, resultErr := testDB.GetInvertedFileFromKey(0)
	if resultErr == nil || len(invFileResult) != 0 {
		t.FailNow()
	}

}

func TestAddGetDatabaseInvertedFileIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &InvertedFileIndexer{}
	err := testDB.Initialize(wd + "/dbTest/InvertedFileIndexer")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}

	testDB.DeleteAllInvertedFileFromKey(0)

	invFile := InvertedFile{pageID: 0}
	invFile.AddWordPositions(0)
	invFile.AddWordPositions(1)

	addErr := testDB.AddKeyToIndexOrUpdate(0, invFile)
	if addErr != nil {
		t.FailNow()
	}

	invFile2 := InvertedFile{pageID: 1}
	invFile2.AddWordPositions(2)
	invFile2.AddWordPositions(3)

	addErr = testDB.AddKeyToIndexOrUpdate(0, invFile2)
	if addErr != nil {
		t.FailNow()
	}

	invFileResult, resultErr := testDB.GetInvertedFileFromKey(0)
	if resultErr != nil {
		t.FailNow()
	}

	if len(invFileResult) != 2 {
		t.Fail()
	}

	if !(invFileResult[0].Same(&invFile) && invFileResult[1].Same(&invFile2)) {
		t.Fail()
	}
}

func TestDeleteDatabaseInvertedFileIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &InvertedFileIndexer{}
	err := testDB.Initialize(wd + "/dbTest/InvertedFileIndexer")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}

	invFile := InvertedFile{pageID: 0}
	invFile.AddWordPositions(0)
	invFile.AddWordPositions(1)

	addErr := testDB.AddKeyToIndexOrUpdate(0, invFile)
	if addErr != nil {
		t.FailNow()
	}

	invFile2 := InvertedFile{pageID: 1}
	invFile2.AddWordPositions(2)
	invFile2.AddWordPositions(3)

	addErr = testDB.AddKeyToIndexOrUpdate(0, invFile2)
	if addErr != nil {
		t.FailNow()
	}

	wordList := make([]uint64, 1)
	wordList[0] = 0

	testDB.DeleteInvertedFileFromWordListAndPage(wordList, 0)

	invFileResult, resultErr := testDB.GetInvertedFileFromKey(0)
	if resultErr != nil {
		t.FailNow()
	}

	if len(invFileResult) > 2 {
		t.Fail()
	}

	if invFileResult[0].Same(&invFile) {
		t.Fail()
	}
}

func TestInitializePagePropetiesIndexer(t *testing.T) {

	wd, _ := os.Getwd()
	testDB := &PagePropetiesIndexer{}
	err := testDB.Initialize(wd + "/dbTest/PagePropetiesIndexer")
	defer testDB.Release()
	if err != nil {
		t.Fail()
	}
}

func TestIteratePagePropetiesIndexer(t *testing.T) {

	wd, _ := os.Getwd()
	testDB := &PagePropetiesIndexer{}
	err := testDB.Initialize(wd + "/dbTest/PagePropetiesIndexer")
	defer testDB.Release()
	if err != nil {
		t.Fail()
	}
	testDB.Iterate()
}

func TestAddGetDatabasePagePropetiesIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &PagePropetiesIndexer{}
	err := testDB.Initialize(wd + "/dbTest/PagePropetiesIndexer")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}
	page := Page{id: 0, title: "Test Page", url: "www.testpage.com", size: 10}
	addErr := testDB.AddKeyToPageProperties(0, page)
	if addErr != nil {
		t.FailNow()
	}

	pageResult, resultErr := testDB.GetPagePropertiesFromKey(0)
	if resultErr != nil {
		t.FailNow()
	}
	if pageResult.GetId() != 0 || pageResult.GetTitle() != "Test Page" || pageResult.GetUrl() != "www.testpage.com" || pageResult.GetSize() != 10 {
		t.Fail()
	}
}

func TestAddGetAllPagePropetiesIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &PagePropetiesIndexer{}
	err := testDB.Initialize(wd + "/dbTest/PagePropetiesIndexer")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}
	page := Page{id: 0, title: "Test Page", url: "www.testpage.com", size: 10}
	page2 := CreatePage(page.GetId(), page.GetTitle(), page.GetUrl(), page.GetSize(), time.Now())

	addErr := testDB.AddKeyToPageProperties(0, page)
	addErr = testDB.AddKeyToPageProperties(1, page2)
	if addErr != nil {
		t.FailNow()
	}
	pageList, allErr := testDB.All()
	if allErr != nil {
		t.FailNow()
	}
	if len(pageList) == 0 {
		t.FailNow()
	}
}

func TestDeleteDatabasePagePropetiesIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &PagePropetiesIndexer{}
	err := testDB.Initialize(wd + "/dbTest/PagePropetiesIndexer")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}
	page := Page{id: 0, title: "Test Page", url: "www.testpage.com", size: 10}
	addErr := testDB.AddKeyToPageProperties(0, page)
	if addErr != nil {
		t.FailNow()
	}

	deleteErr := testDB.DeletePagePropertiesFromKey(0)
	if deleteErr != nil {
		t.FailNow()
	}
	_, resultErr := testDB.GetPagePropertiesFromKey(0)
	if resultErr == nil {
		t.FailNow()
	}
}

func TestInitializeReverseMappingIndexer(t *testing.T) {

	wd, _ := os.Getwd()
	testDB := &ReverseMappingIndexer{}
	err := testDB.Initialize(wd + "/dbTest/ReverseMappingIndex")
	defer testDB.Release()
	if err != nil {
		t.Fail()
	}
}

func TestIterateReverseMappingIndexer(t *testing.T) {

	wd, _ := os.Getwd()
	testDB := &ReverseMappingIndexer{}
	err := testDB.Initialize(wd + "/dbTest/ReverseMappingIndex")
	defer testDB.Release()
	if err != nil {
		t.Fail()
	}
	testDB.Iterate()
}

func TestAddGetDatabaseReverseMappingIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &ReverseMappingIndexer{}
	err := testDB.Initialize(wd + "/dbTest/ReverseMappingIndex")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}

	addErr := testDB.AddKeyToIndex(0, "test")
	if addErr != nil {
		t.FailNow()
	}

	idResult, resultErr := testDB.GetValueFromKey(0)
	if resultErr != nil {
		t.FailNow()
	}

	if idResult != "test" {
		t.Fail()
	}
}

func TestDeleteDatabaseReverseMappingIndexer(t *testing.T) {
	wd, _ := os.Getwd()
	testDB := &ReverseMappingIndexer{}
	err := testDB.Initialize(wd + "/dbTest/ReverseMappingIndex")
	defer testDB.Release()
	if err != nil {
		t.FailNow()
	}

	idList := make([]uint64, 1)
	idList[0] = 0

	addErr := testDB.AddKeyToIndex(0, "test")
	if addErr != nil {
		t.FailNow()
	}

	deleteErr := testDB.DeleteKeyValuePair(0)
	if deleteErr != nil {
		t.FailNow()
	}

	_, resultErr := testDB.GetValueFromKey(0)
	if resultErr == nil {
		t.FailNow()
	}

}
