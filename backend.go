package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/davi1972/comp4321-search-engine/boolsearch"
	"github.com/davi1972/comp4321-search-engine/phrasalSearch"
	"github.com/davi1972/comp4321-search-engine/vsm"
	"github.com/dgraph-io/badger"

	"github.com/gorilla/mux"

	Indexer "github.com/davi1972/comp4321-search-engine/indexer"
)

type server struct {
	documentIndexer                   *Indexer.MappingIndexer
	wordIndexer                       *Indexer.MappingIndexer
	reverseDocumentIndexer            *Indexer.ReverseMappingIndexer
	reverseWordIndexer                *Indexer.ReverseMappingIndexer
	pagePropertiesIndexer             *Indexer.PagePropetiesIndexer
	titleInvertedIndexer              *Indexer.InvertedFileIndexer
	contentInvertedIndexer            *Indexer.InvertedFileIndexer
	documentWordForwardIndexer        *Indexer.DocumentWordForwardIndexer
	titleWordForwardIndexer           *Indexer.DocumentWordForwardIndexer
	parentChildDocumentForwardIndexer *Indexer.ForwardIndexer
	childParentDocumentForwardIndexer *Indexer.ForwardIndexer
	wordCountContentIndexer           *Indexer.PageRankIndexer
	pageRankIndexer                   *Indexer.PageRankIndexer
	router                            *mux.Router
	vsm                               *vsm.VSM
	bs                                *boolsearch.BoolSearch
	pls                               *phrasalSearch.PhrasalSearch
}

type Edge struct {
	From int `json:"source"`
	To   int `json:"target"`
}

type EdgeString struct {
	From string `json:"source"`
	To   string `json:"target"`
}

type Node struct {
	Name string `json:"id"`
}

type GraphResponse struct {
	Nodes       []Node `json:"nodes"`
	Edges       []Edge `json:"links"`
	EdgesString []EdgeString
}

type WordListResponse struct {
	WordList []string `json:"words"`
}

type WordFrequencyString struct {
	Word      string `json:"word"`
	Frequency uint64 `json:"frequency"`
}

type QueryResponse struct {
	VSMScore         float64               `json:"vsmscore"`
	PageRankScore    float64               `json:"pagerankscore"`
	Score            float64               `json:"score"`
	Title            string                `json:"title"`
	URL              string                `json:"url"`
	LastModifiedDate time.Time             `json:"last_modified"`
	KeyWord          []WordFrequencyString `json:"keywords"`
	ParentList       []string              `json:"parent_urls"`
	ChildList        []string              `json:"child_urls"`
}

type QueryResponses []QueryResponse

func (s QueryResponses) Len() int {
	return len(s)
}

func (s QueryResponses) Less(i, j int) bool {
	return s[i].Score > s[j].Score
}

func (s QueryResponses) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type QueryListResponse struct {
	List QueryResponses `json:"documents"`
}

// S ...
var S server
var maxDepth = 2
var prWeight = 0.8

func main() {
	S.Initialize()
	S.routes()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		S.Release()
		os.Exit(1)
	}()
	http.ListenAndServe("localhost:8000", S.router)
}

func (s *server) Initialize() {
	wd, _ := os.Getwd()
	s.documentIndexer = &Indexer.MappingIndexer{}
	docErr := s.documentIndexer.Initialize(wd + "/db/documentIndex")
	if docErr != nil {
		fmt.Printf("error when initializing document indexer: %s\n", docErr)
	}

	s.reverseDocumentIndexer = &Indexer.ReverseMappingIndexer{}
	reverseDocumentIndexerErr := s.reverseDocumentIndexer.Initialize(wd + "/db/reverseDocumentIndexer")
	if reverseDocumentIndexerErr != nil {
		fmt.Printf("error when initializing reverse document indexer: %s\n", reverseDocumentIndexerErr)
	}

	s.wordIndexer = &Indexer.MappingIndexer{}
	wordErr := s.wordIndexer.Initialize(wd + "/db/wordIndex")
	if wordErr != nil {
		fmt.Printf("error when initializing word indexer: %s\n", wordErr)
	}

	s.reverseWordIndexer = &Indexer.ReverseMappingIndexer{}
	reverseWordindexerErr := s.reverseWordIndexer.Initialize(wd + "/db/reverseWordIndexer")
	if reverseWordindexerErr != nil {
		fmt.Printf("error when initializing reverse word indexer: %s\n", reverseWordindexerErr)
	}

	s.pagePropertiesIndexer = &Indexer.PagePropetiesIndexer{}
	pagePropertiesErr := s.pagePropertiesIndexer.Initialize(wd + "/db/pagePropertiesIndex")
	if pagePropertiesErr != nil {
		fmt.Printf("error when initializing page properties: %s\n", pagePropertiesErr)
	}

	s.titleInvertedIndexer = &Indexer.InvertedFileIndexer{}
	titleInvertedErr := s.titleInvertedIndexer.Initialize(wd + "/db/titleInvertedIndex")
	if titleInvertedErr != nil {
		fmt.Printf("error when initializing page properties: %s\n", titleInvertedErr)
	}

	s.contentInvertedIndexer = &Indexer.InvertedFileIndexer{}
	contentInvertedErr := s.contentInvertedIndexer.Initialize(wd + "/db/contentInvertedIndex")
	if contentInvertedErr != nil {
		fmt.Printf("error when initializing page properties: %s\n", contentInvertedErr)
	}

	s.documentWordForwardIndexer = &Indexer.DocumentWordForwardIndexer{}
	documentWordForwardIndexerErr := s.documentWordForwardIndexer.Initialize(wd + "/db/documentWordForwardIndex")
	if documentWordForwardIndexerErr != nil {
		fmt.Printf("error when initializing document -> word forward Indexer: %s\n", documentWordForwardIndexerErr)
	}

	s.parentChildDocumentForwardIndexer = &Indexer.ForwardIndexer{}
	parentChildDocumentForwardIndexerErr := s.parentChildDocumentForwardIndexer.Initialize(wd + "/db/parentChildDocumentForwardIndex")
	if parentChildDocumentForwardIndexerErr != nil {
		fmt.Printf("error when initializing parentDocument -> childDocument forward Indexer: %s\n", parentChildDocumentForwardIndexerErr)
	}

	s.childParentDocumentForwardIndexer = &Indexer.ForwardIndexer{}
	childParentDocumentForwardIndexerErr := s.childParentDocumentForwardIndexer.Initialize(wd + "/db/childParentDocumentForwardIndex")
	if childParentDocumentForwardIndexerErr != nil {
		fmt.Printf("error when initializing childDocument -> parentDocument forward Indexer: %s\n", childParentDocumentForwardIndexerErr)
	}

	s.titleWordForwardIndexer = &Indexer.DocumentWordForwardIndexer{}
	titleWordForwardIndexerErr := s.titleWordForwardIndexer.Initialize(wd + "/db/titleWordForwardIndex")
	if titleWordForwardIndexerErr != nil {
		fmt.Printf("error when initializing document -> word forward Indexer: %s\n", titleWordForwardIndexerErr)
	}

	s.pageRankIndexer = &Indexer.PageRankIndexer{}
	pageRankIndexerErr := s.pageRankIndexer.Initialize(wd + "/db/pageRankIndex")
	if pageRankIndexerErr != nil {
		fmt.Printf("error when initializing page rank indexer: %s\n", pageRankIndexerErr)
	}

	s.router = mux.NewRouter()
	s.vsm = &vsm.VSM{
		DocumentIndexer:                   s.documentIndexer,
		WordIndexer:                       s.wordIndexer,
		ReverseDocumentIndexer:            s.reverseDocumentIndexer,
		ReverseWordIndexer:                s.reverseWordIndexer,
		PagePropertiesIndexer:             s.pagePropertiesIndexer,
		TitleInvertedIndexer:              s.titleInvertedIndexer,
		ContentInvertedIndexer:            s.contentInvertedIndexer,
		DocumentWordForwardIndexer:        s.documentWordForwardIndexer,
		ParentChildDocumentForwardIndexer: s.parentChildDocumentForwardIndexer,
		ChildParentDocumentForwardIndexer: s.childParentDocumentForwardIndexer,
		TitleWordForwardIndexer:           s.titleWordForwardIndexer,
	}

	s.bs = &boolsearch.BoolSearch{
		ContentInvertedIndexer: s.contentInvertedIndexer,
		Vsm:                    s.vsm,
	}

	s.pls = &phrasalSearch.PhrasalSearch{
		TitleInvertedIndexer:    s.titleInvertedIndexer,
		ContentInvertedIndexer:  s.contentInvertedIndexer,
		TitleWordForwardIndexer: s.titleWordForwardIndexer,
		V:                       s.vsm,
		Bs:                      s.bs,
	}

}

func (s *server) Release() {
	s.documentIndexer.Release()
	s.reverseDocumentIndexer.Release()
	s.wordIndexer.Release()
	s.reverseWordIndexer.Release()
	s.pagePropertiesIndexer.Release()
	s.titleInvertedIndexer.Release()
	s.contentInvertedIndexer.Release()
	s.documentWordForwardIndexer.Release()
	s.parentChildDocumentForwardIndexer.Release()
	s.childParentDocumentForwardIndexer.Release()
	s.pageRankIndexer.Release()
	s.titleWordForwardIndexer.Release()
}

func (g *GraphResponse) AppendNodesAndEdgesStringFromIDList(docIDs []uint64) ([]uint64, error) {
	resultIDs := []uint64{}
	for _, docID := range docIDs {
		curStr, curErr := S.reverseDocumentIndexer.GetValueFromKey(docID)
		if curErr != nil {
			continue
		}
		idList, _ := S.parentChildDocumentForwardIndexer.GetIdListFromKey(uint64(docID))
		for _, i := range idList {
			str, valErr := S.reverseDocumentIndexer.GetValueFromKey(i)
			if valErr == badger.ErrKeyNotFound {
				continue
			} else if valErr != nil {
				return nil, valErr
			}

			g.EdgesString = append(g.EdgesString, EdgeString{From: curStr, To: str})
			exist := false
			for _, node := range g.Nodes {
				if node.Name == str {
					exist = true
				}
			}
			if !exist {
				g.Nodes = append(g.Nodes, Node{Name: str})
			}

		}
		resultIDs = append(resultIDs, idList...)
	}
	return resultIDs, nil
}

// CreateEdgesID ...
func (g *GraphResponse) CreateEdgesID() {
	idMap := make(map[string]int)
	for i, val := range g.Nodes {
		if _, ok := idMap[val.Name]; !ok {
			idMap[val.Name] = i
		}
	}
	for _, val := range g.EdgesString {
		g.Edges = append(g.Edges, Edge{From: idMap[val.From], To: idMap[val.To]})
	}
}

func graphHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, convertErr := strconv.Atoi(vars["documentID"])
	if convertErr != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Invalid parameter value! Details: " + convertErr.Error()))
	}
	resp := &GraphResponse{}
	curIDList := []uint64{}
	var iterErr error
	// Append first id to curIDList
	curIDList = append(curIDList, uint64(id))
	for iterations := 0; iterations < maxDepth; iterations++ {
		curIDList, iterErr = resp.AppendNodesAndEdgesStringFromIDList(curIDList)
		if iterErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 - Internal Server Error! Details: " + iterErr.Error()))
			break
		}
	}
	resp.CreateEdgesID()
	w.Header().Set("Access-Control-Allow-Origin", "*")
	jsonResult, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}

func wordListHandler(w http.ResponseWriter, r *http.Request) {
	resp := &WordListResponse{}
	resp.WordList = S.wordIndexer.AllValue()
	jsonResult, _ := json.Marshal(resp)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
}

func queryHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	query := vars["queryString"]

	// Extract phrases first before doing everything else
	regex, _ := regexp.Compile(`("([^"]|"")*")`)
	phraseList := regex.FindAllString(query, -1)
	boostedDocsIDList := make(map[uint64]int)
	for _, phrase := range phraseList {
		splitPhrase := strings.Split(strings.Trim(phrase, "\""), " ")
		for _, doc := range S.pls.GetPhraseDocuments(splitPhrase) {
			boostedDocsIDList[doc]++
		}
	}

	resp := &QueryListResponse{}

	responses := QueryResponses{}

	start := time.Now()
	cosScore, err := S.vsm.ComputeCosineScore(query)
	elapsed := time.Since(start)
	log.Printf("Cosine took %s", elapsed)
	start = time.Now()

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 - Invalid parameter value! Details: " + err.Error()))
	}
	for i, score := range cosScore {
		if score == 0 {
			continue
		}

		pageRankScore, err := S.pageRankIndexer.GetValueFromKey(i)

		if math.IsInf(pageRankScore, 1) {
			pageRankScore = 0
		}

		if err != nil {
			fmt.Println("error retrieving page rank value", err)
		}

		doc := &QueryResponse{}
		doc.PageRankScore = pageRankScore

		doc.VSMScore = score
		doc.Score = prWeight*pageRankScore + (1-prWeight)*score
		// add boost to phrases!
		if _, ok := boostedDocsIDList[i]; ok {
			fmt.Println("Applying boost as phrase search")
			doc.Score *= 1.5
		}
		pageProps, _ := S.pagePropertiesIndexer.GetPagePropertiesFromKey(i)
		doc.Title = pageProps.GetTitle()
		doc.URL = pageProps.GetUrl()
		doc.LastModifiedDate = pageProps.GetDate()
		childList, _ := S.parentChildDocumentForwardIndexer.GetIdListFromKey(i)
		for _, childID := range childList {
			str, err := S.reverseDocumentIndexer.GetValueFromKey(childID)
			if err == nil {
				doc.ChildList = append(doc.ChildList, str)
			}
		}
		parentList, _ := S.childParentDocumentForwardIndexer.GetIdListFromKey(i)
		for _, parentID := range parentList {
			str, err := S.reverseDocumentIndexer.GetValueFromKey(parentID)
			if err == nil {
				doc.ParentList = append(doc.ParentList, str)
			}
		}

		wordFreq, _ := S.documentWordForwardIndexer.GetWordFrequencyListFromKey(i)
		sort.Sort(Indexer.WordFrequencySorter(wordFreq))
		for _, wordF := range wordFreq[:5] {
			wordStr, wordErr := S.reverseWordIndexer.GetValueFromKey(wordF.GetID())
			if wordErr == nil {
				doc.KeyWord = append(doc.KeyWord, WordFrequencyString{Word: wordStr, Frequency: wordF.GetFrequency()})
			}
		}

		responses = append(responses, *doc)
	}
	sort.Sort(responses)
	resp.List = responses
	jsonResult, jsonErr := json.Marshal(resp)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal Server Error! Details: " + jsonErr.Error()))
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResult)
	elapsed = time.Since(start)
	log.Printf("Forming response took %s", elapsed)
}

func (s *server) routes() {
	s.router.HandleFunc("/graph/{documentID}", graphHandler)
	s.router.HandleFunc("/wordList", wordListHandler)
	s.router.HandleFunc("/query/{queryString}", queryHandler)
}
