package models

import "database/sql"

type Company struct {
	ID        uint64 `json:"id"`
	Name      string `json:"name"`
	TickerUrl string `json:"ticker_url"`
}

type Document struct {
	ID            uint64
	CompanyName   string
	DocType       string
	Url           sql.NullString
	PublishedDate sql.NullString
	IsTrained     uint8
	ModifiedAt    string
	Content       []byte
}

type Sector struct {
	ID     int64  `json:"id"`
	Sector string `json:"sector"`
}
type ExchangeFilingDoc struct {
	ID          int64  `json:"id"`
	CompanyID   int64  `json:"company_id"`
	PublishedIn string `json:"published_in"`
	DocURL      string `json:"doc_url"`
}
