package tfidf

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/dchest/stemmer/porter2"
	"github.com/tickrlytics/tickerlytics-backend/configs"
	"github.com/tickrlytics/tickerlytics-backend/controllers/tfidf/tokenize"
	"github.com/tickrlytics/tickerlytics-backend/models"
)

const (
	K                    float64 = 1.2
	B                    float64 = 0.75
	MIN_WORD_LEN_ALLOWED         = 2
	MAX_WORD_LEN_ALLOWED         = 15
	SEPRATORS                    = `,.[]{}()<>/|\<>:;?`
)

var (
	tfIdfInstance    *tfidf
	lastModifiedTime time.Time
	lastReadTime     time.Time
	lock             = &sync.Mutex{}
)

// tfidf tfidf model
type tfidf struct {
	Id           string            `json:"id"`
	DocIndex     map[string]uint64 `json:"doc_index"`     // train document index
	IndexDocName map[uint64]string `json:"index_docname"` // train document name mapped with index
	//termFreqs     []map[string]uint64       // terms frequency for each train document
	DocsTermsFreq map[uint64]map[string]uint64 `json:"docs_terms_freq"` // terms frequency for each train document
	TermDocsCount map[string]uint64            `json:"term_docs_count"` // number of document in which term t appears in train data
	TermDocs      map[string]map[uint64]bool   `json:"term_docs"`       //list of doc in which term t present.
	//docsTermLen   []uint64                  //len of each doc in words.
	DocsTermLenMap    map[uint64]uint64      `json:"docs_term_len"` //len of each doc in words.
	DocsPublishedTime map[uint64]string      `json:"docs_published_time"`
	DocsModifiedTime  map[uint64]string      `json:"docs_modified_time"`
	DocsUrl           map[uint64]string      `json:"docs_url"`
	DocsCompany       map[uint64]string      `json:"docs_company"`
	TermsLen          uint64                 `json:"terms_len"`  //total terms in corpus
	TotalDocs         uint64                 `json:"total_docs"` // number of documents in train data
	ModifiedTime      string                 `json:"modified_time"`
	stopWords         map[string]interface{} // words to be remove.
	tokenizer         tokenize.Tokenizer
}

// New new model with default
func New() (*tfidf, error) {

	if tfIdfInstance != nil && lastModifiedTime.Before(lastReadTime) {
		return tfIdfInstance, nil
	}
	lock.Lock()
	defer lock.Unlock()
	if tfIdfInstance != nil && lastModifiedTime.Before(lastReadTime) {
		return tfIdfInstance, nil
	}
	tfidfID := configs.GetEnvWithKey("AZ_COSMOS_DB_TFIDF_ITEM_ID", "")
	err := readTFIDF(tfidfID)
	if err != nil {
		err, isExportedErr := err.(*azcore.ResponseError)
		if isExportedErr && err.StatusCode == http.StatusNotFound {
			tfIdfInstance = new(tfidf)
			tfIdfInstance.Id = tfidfID
		} else {
			return nil, err
		}
	}
	lastReadTime = time.Now()
	tfIdfInstance.tokenizer = &tokenize.EnTokenizer{
		Seprators: SEPRATORS,
	}
	if tfIdfInstance.DocIndex == nil {
		tfIdfInstance.DocIndex = make(map[string]uint64)
	}
	if tfIdfInstance.IndexDocName == nil {
		tfIdfInstance.IndexDocName = make(map[uint64]string)
	}
	if tfIdfInstance.DocsTermsFreq == nil {
		tfIdfInstance.DocsTermsFreq = make(map[uint64]map[string]uint64)
	}
	if tfIdfInstance.TermDocsCount == nil {
		tfIdfInstance.TermDocsCount = make(map[string]uint64)
	}
	if tfIdfInstance.TermDocs == nil {
		tfIdfInstance.TermDocs = map[string]map[uint64]bool{}
	}

	if tfIdfInstance.DocsTermLenMap == nil {
		tfIdfInstance.DocsTermLenMap = make(map[uint64]uint64)
	}
	if tfIdfInstance.DocsPublishedTime == nil {
		tfIdfInstance.DocsPublishedTime = make(map[uint64]string)
	}
	if tfIdfInstance.DocsModifiedTime == nil {
		tfIdfInstance.DocsModifiedTime = make(map[uint64]string)
	}
	if tfIdfInstance.DocsUrl == nil {
		tfIdfInstance.DocsUrl = make(map[uint64]string)
	}
	if tfIdfInstance.DocsCompany == nil {
		tfIdfInstance.DocsCompany = make(map[uint64]string)
	}

	//We have to get below values from cosmosDB
	// tfIdfInstance = &tfidf{
	// 	Id:           tfidfID,
	// tfIdfInstance.DocIndex : make(map[string]uint64)
	// 	IndexDocName: make(map[uint64]string),
	// 	//termFreqs:     make([]map[string]uint64, 0),
	// 	DocsTermsFreq: make(map[uint64]map[string]uint64),
	// 	TermDocs:      make(map[string]map[uint64]bool),
	// 	TermDocsCount: make(map[string]uint64),
	// 	//docsTermLen:   make([]uint64, 0),
	// 	DocsTermLenMap:    make(map[uint64]uint64),
	// 	DocsPublishedTime: make(map[uint64]string),
	// 	DocsUrl:           make(map[uint64]string),
	// 	DocsCompany:       make(map[uint64]string),
	// 	DocsModifiedTime:  make(map[uint64]string),
	// 	TermsLen:          0,
	// 	TotalDocs:         0,
	// 	tokenizer:         &tokenize.EnTokenizer{},
	// }
	stopwords := []string{`i`, `me`, `my`, `myself`, `we`, `our`, `ours`, `ourselves`, `you`, "you`re", "you`ve", "you`ll", "you`d", `your`, `yours`, `yourself`, `yourselves`, `he`, `him`, `his`, `himself`, `she`, "she`s", `her`, `hers`, `herself`, `it`, "it`s", `its`, `itself`, `they`, `them`, `their`, `theirs`, `themselves`, `what`, `which`, `who`, `whom`, `this`, `that`, "that`ll", `these`, `those`, `am`, `is`, `are`, `was`, `were`, `be`, `been`, `being`, `have`, `has`, `had`, `having`, `do`, `does`, `did`, `doing`, `a`, `an`, `the`, `and`, `but`, `if`, `or`, `because`, `as`, `until`, `while`, `of`, `at`, `by`, `for`, `with`, `about`, `against`, `between`, `into`, `through`, `during`, `before`, `after`, `above`, `below`, `to`, `from`, `up`, `down`, `in`, `out`, `on`, `off`, `over`, `under`, `again`, `further`, `then`, `once`, `here`, `there`, `when`, `where`, `why`, `how`, `all`, `any`, `both`, `each`, `few`, `more`, `most`, `other`, `some`, `such`, `no`, `nor`, `not`, `only`, `own`, `same`, `so`, `than`, `too`, `very`, `s`, `t`, `can`, `will`, `just`, `don`, "don`t", `should`, "should`ve", `now`, `d`, `ll`, `m`, `o`, `re`, `ve`, `y`, `ain`, `aren`, "aren`t", `couldn`, "couldn`t", `didn`, "didn`t", `doesn`, "doesn`t", `hadn`, "hadn`t", `hasn`, "hasn`t", `haven`, "haven`t", `isn`, "isn`t", `ma`, `mightn`, "mightn`t", `mustn`, "mustn`t", `needn`, "needn`t", `shan`, "shan`t", `shouldn`, "shouldn`t", `wasn`, "wasn`t", `weren`, "weren`t", `won`, "won`t", `wouldn`, "wouldn`t", "https", `–`}
	tfIdfInstance.AddStopWords(stopwords...)
	return tfIdfInstance, nil
}

// AddStopWords add stop words to be remove
func (tfIDF *tfidf) AddStopWords(words ...string) {
	if tfIDF.stopWords == nil {
		tfIDF.stopWords = make(map[string]interface{})
	}

	for _, word := range words {
		tfIDF.stopWords[word] = nil
	}
}

func (tfIDF *tfidf) TrainDoc(doc *models.Document) {
	docHash := hash(doc.Content)
	docPos := tfIDF.docHashPos(docHash)
	if docPos >= 1 {
		doc.IsTrained = configs.STATUS_DOC_ALEREADY_TRAINED
		return
	}
	tokens := tfIDF.tokenizer.Tokens(string(doc.Content))
	termFreq := tfIDF.termFreq(tokens)
	//not required to train doc as it does not contain useful information.
	if len(termFreq) == 0 {
		doc.IsTrained = configs.STATUS_DOC_TRAINING_NOT_REQUIRED
		return
	}
	lock.Lock()
	defer lock.Unlock()
	tfIDF.TotalDocs += 1
	var docTokenCount uint64 = 0
	for term, freq := range termFreq {
		termDocSet, isFound := tfIDF.TermDocs[term]
		if !isFound {
			termDocSet = make(map[uint64]bool)
		}
		docTokenCount += freq
		// termDocSet[tfIDF.n] = true
		termDocSet[doc.ID] = true
		tfIDF.TermDocs[term] = termDocSet
		tfIDF.TermDocsCount[term] += 1
	}

	// tfIDF.termFreqsMap[tfIDF.n] = termFreq
	tfIDF.DocsTermsFreq[doc.ID] = termFreq
	// tfIDF.docsTermLenMap[tfIDF.n] = docTokenCount
	tfIDF.DocsTermLenMap[doc.ID] = docTokenCount
	tfIDF.TermsLen = tfIDF.TermsLen + docTokenCount
	// tfIDF.docIndex[docHash] = tfIDF.n
	tfIDF.DocIndex[docHash] = doc.ID
	// tfIDF.indexDocName[tfIDF.n] = doc.DocType
	tfIDF.IndexDocName[doc.ID] = doc.DocType
	// tfIDF.docsCompany[tfIDF.n] = doc.CompanyName
	tfIDF.DocsCompany[doc.ID] = doc.CompanyName
	// tfIDF.docsUrl[tfIDF.n] = doc.Url.String
	tfIDF.DocsUrl[doc.ID] = doc.Url.String
	// tfIDF.docsPublishedDate[tfIDF.n] = doc.PublishedDate.String
	tfIDF.DocsPublishedTime[doc.ID] = doc.PublishedDate.String
	// tfIDF.docsModifiedTime[tfIDF.n] = time.Now().Format("2006-01-02")
	modifiedTime := time.Now()
	tfIDF.ModifiedTime = modifiedTime.Format("2006-01-02 15:04:05")
	tfIDF.DocsModifiedTime[doc.ID] = tfIDF.ModifiedTime
	doc.IsTrained = configs.STATUS_DOC_TRAINED
	go upsertTFDIF(modifiedTime, doc.ID)
}
func upsertTFDIF(modifiedTime time.Time, docId uint64) {
	if modifiedTime.Before(lastModifiedTime) {
		configs.Logger.Infof("trained data for doc_id %d is already inserted", docId)
		return
	}
	data, err := json.Marshal(tfIdfInstance)
	if err != nil {
		configs.Logger.Errorf("error in marshalling tfidf: %s", err.Error())
		return
	}
	lock.Lock()
	defer lock.Unlock()
	if modifiedTime.Before(lastModifiedTime) {
		configs.Logger.Infof("trained data for doc_id %d is already inserted", docId)
		return
	}
	cosmos, err := models.GetCosmoDBInstance()
	if err != nil {
		configs.Logger.Errorf("error in getting cosmos instance: %s", err.Error())
		return
	}
	err = cosmos.UpsertItem(data, tfIdfInstance.Id, configs.GetEnvWithKey("AZ_COSMOS_DB_DATABASE_ID", ""), configs.GetEnvWithKey("AZ_COSMOS_DB_CONTAINER_ID", ""))
	if err != nil {
		configs.Logger.Errorf("error in upserting tfidf doc with doc_id %d into the cosmos: %s", docId, err.Error())
		return
	}
	configs.Logger.Infof("trained doc inserted into cosmos: modified_time: %s", modifiedTime.Format("2006-01-02 15:04:05"))
	lastModifiedTime = modifiedTime
}
func readTFIDF(tfidfID string) error {
	cosmos, err := models.GetCosmoDBInstance()
	if err != nil {
		configs.Logger.Errorf("error in getting cosmos instance: %s", err.Error())
		return err
	}
	res, err := cosmos.GetItem(tfidfID, configs.GetEnvWithKey("AZ_COSMOS_DB_DATABASE_ID", ""), configs.GetEnvWithKey("AZ_COSMOS_DB_CONTAINER_ID", ""))
	if err != nil {
		return err
	}
	tfIdfInstance = new(tfidf)
	err = json.Unmarshal(res, tfIdfInstance)
	if err != nil {
		return err
	}
	return nil
}

// DocName return document name
func (tfIDF *tfidf) DocName(docPos uint64) string {
	return tfIDF.IndexDocName[docPos]
}

// DocUrl return document url
func (tfIDF *tfidf) DocUrl(docPos uint64) string {
	return tfIDF.DocsUrl[docPos]
}

// DocPublishedDate return published date of document
func (tfIDF *tfidf) DocPublishedDate(docPos uint64) string {
	return tfIDF.DocsPublishedTime[docPos]
}

// termFreq calculate term-freq of each term in document.
func (tfIDF *tfidf) termFreq(tokens []string) (m map[string]uint64) {
	m = make(map[string]uint64)
	for _, term := range tokens {
		term = cleanWord(term)
		if _, ok := tfIDF.stopWords[term]; ok {
			continue
		}
		term = porter2.Stemmer.Stem(term)
		if len(term) <= MIN_WORD_LEN_ALLOWED || len(term) > MAX_WORD_LEN_ALLOWED {
			continue
		} else if _, ok := tfIDF.stopWords[term]; ok {
			continue
		}
		freq, isFound := m[term]
		if !isFound {
			freq = 1
			m[term] = freq
		} else {
			m[term] = freq + 1
		}
	}
	return
}

// docHashPos return the position of doc in corpus.
func (tfIDF *tfidf) docHashPos(hash string) uint64 {
	if pos, ok := tfIDF.DocIndex[hash]; ok {
		return pos
	}

	return 0
}

// docPos return the position of doc in corpus.
func (tfIDF *tfidf) docPos(doc string) uint64 {
	return tfIDF.docHashPos(hash([]byte(doc)))
}

// hash return hash of the doc content.
func hash(text []byte) string {
	h := md5.New()
	h.Write(text)
	return hex.EncodeToString(h.Sum(nil))
}

// findTfIdf calculate tf-idf.
func findTfIdf(termFreq, docTerms, termDocs, N int) float64 {
	tf := float64(termFreq) / float64(docTerms)
	idf := math.Log(float64(1+N) / (float64(1 + termDocs)))
	// idf := (float64(1+N) / float64(1+termDocs))
	return tf * idf
}

// BM25 calculate BM25 score of the term.
func (tfIDF *tfidf) BM25(docPos uint64, term string) float64 {
	// termFreq := tfIDF.termFreqs[docPos-1][term]
	termFreq := tfIDF.DocsTermsFreq[docPos][term]
	termDocCount := tfIDF.TermDocsCount[term]
	docLen := tfIDF.DocsTermLenMap[docPos]
	docsLen := tfIDF.TermsLen
	IDF := math.Log(1 + ((float64(tfIDF.TotalDocs) - float64(termDocCount) + 0.5) / (float64(termDocCount) + 0.5)))

	avgDocsLen := float64(docsLen) / float64(tfIDF.TotalDocs)
	numerator := float64(termFreq) * (K + 1)
	deno := float64(termFreq) + K*(1-B+B*(float64(docLen)/avgDocsLen))
	bm25 := IDF * (numerator / deno)
	return bm25

}
