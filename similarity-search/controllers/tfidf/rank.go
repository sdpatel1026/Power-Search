package tfidf

import "github.com/dchest/stemmer/porter2"

type Doc struct {
	Rank                float64 `json:"rank"`
	ReleventResultCount uint64  `json:"relevent_results"`
	DocType             string  `json:"doc_type"`
	Url                 string  `json:"url"`
	PublishedDate       string  `json:"published_date"`
}

// CompanyDocs maps company-name with Docs.
type CompanyDocs map[string]Docs

// Docs maps document-id(document position) with Doc.
type Docs map[uint64]Doc

// FindCompanyDocs calculate BM25 rank of the document for the given query
func (tfIDF *tfidf) FindCompanyDocs(query string) CompanyDocs {
	tokens := tfIDF.tokenizer.Tokens(query)
	companyDocs := make(CompanyDocs)
	allTokensPresentDocs := make(map[uint64]uint8)
	for _, token := range tokens {
		if _, ok := tfIDF.stopWords[token]; ok {
			continue
		}
		token = cleanWord(token)
		if len(token) < 2 {
			continue
		}
		token = porter2.Stemmer.Stem(token)
		docList, isFound := tfIDF.TermDocs[token]
		if isFound {

			for docPos := range docList {
				company := tfIDF.DocsCompany[docPos]
				companyDoc, isFound := companyDocs[company]
				if !isFound {
					companyDoc = make(Docs)
					companyDocs[company] = companyDoc
				}
				val, isPresent := allTokensPresentDocs[docPos]
				if !isPresent {
					allTokensPresentDocs[docPos] = 0
				} else if val == 1 {
					continue
				}
				_, isTokenPresent := tfIDF.DocsTermsFreq[docPos][token]
				if !isTokenPresent {
					allTokensPresentDocs[docPos] = 1
					delete(companyDoc, docPos)
				}
				bm25 := tfIDF.BM25(docPos, token)
				doc := companyDoc[docPos]
				doc.DocType = tfIDF.DocName(docPos)
				doc.PublishedDate = tfIDF.DocPublishedDate(docPos)
				doc.Url = tfIDF.DocUrl(docPos)
				doc.Rank = doc.Rank + bm25
				doc.ReleventResultCount = doc.ReleventResultCount + tfIDF.DocsTermsFreq[docPos][token]
				companyDoc[docPos] = doc
			}

		}
	}
	return companyDocs
}
