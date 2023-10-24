package models

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tickrlytics/tickerlytics-backend/configs"
)

type MySQL struct {
	db *sql.DB
}

var mySQL *MySQL

func connectMySql(user, password, host, db string) (*sql.DB, error) {
	rootCertPool := x509.NewCertPool()
	pem, _ := ioutil.ReadFile("DigiCertGlobalRootCA.crt.pem")
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return nil, errors.New("Failed to append PEM.")
	}
	mysql.RegisterTLSConfig("custom", &tls.Config{RootCAs: rootCertPool})
	ConnectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s?allowNativePasswords=true&tls=custom", user, password, host, db)
	return sql.Open("mysql", ConnectionString)
}

func GetMySql() *MySQL {
	return mySQL
}

func (mySQL *MySQL) InsertDoc(compID int, docType, url, publishedDate string) error {
	insertQuery := "INSERT into documents(company_id,doc_type,url,published_date) VALUES(?,?,?,?)"
	_, err := mySQL.db.Exec(insertQuery, compID, docType, url, publishedDate)
	return err
}

func (mySQL *MySQL) GetUnTrainedDocs(requestID string, offset, limit int) ([]Document, error) {
	selectQuery := `SELECT D.id,D.doc_type,D.url, D.published_date, D.is_trained, D.modified_time,C.company_name FROM companies AS C INNER JOIN documents AS D ON D.company_id=C.id WHERE D.is_trained=? LIMIT ?,?`
	var documents []Document
	rows, err := mySQL.db.Query(selectQuery, configs.STATUS_DOC_NOT_TRAINED, offset, limit)
	if err != nil {
		return documents, err
	}
	defer rows.Close()
	for rows.Next() {
		var doc Document
		err = rows.Scan(&doc.ID, &doc.DocType, &doc.Url, &doc.PublishedDate, &doc.IsTrained, &doc.ModifiedAt, &doc.CompanyName)
		if err != nil {
			configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
			continue
		}
		documents = append(documents, doc)
	}
	return documents, nil
}

func (mySQL *MySQL) UpdateDocTrainStatus(ID uint64, status uint8) error {
	updateQuery := "UPDATE documents SET is_trained=? WHERE id=?"
	_, err := mySQL.db.Exec(updateQuery, status, ID)
	if err != nil {
		return err
	}
	return nil
}
func (mySQL *MySQL) GetSectors() ([]Sector, error) {
	selectQuery := "SELECT id,sector FROM sectors"
	raws, err := mySQL.db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer raws.Close()
	var sectors []Sector
	for raws.Next() {
		var sector Sector
		err = raws.Scan(&sector.ID, &sector.Sector)
		if err != nil {
			configs.Logger.Errorf("error in scanning raws while retrieving sectors:%s", err.Error())
			continue
		}
		sectors = append(sectors, sector)
	}
	return sectors, nil
}

func (mySQL *MySQL) GetCompanies() ([]Company, error) {
	selectQuery := "SELECT id,company_name,ticker_url FROM companies"
	raws, err := mySQL.db.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer raws.Close()
	var companies []Company
	for raws.Next() {
		var company Company
		var tickerURL sql.NullString
		err = raws.Scan(&company.ID, &company.Name, &tickerURL)
		if err != nil {
			configs.Logger.Errorf("error in scanning raws while retrieving companies:%s", err.Error())
			continue
		}
		company.TickerUrl = tickerURL.String
		companies = append(companies, company)
	}
	return companies, nil
}
func (mySQL *MySQL) GetSectorCompanies(id int64) (map[uint64]Company, error) {
	selectQuery := "SELECT C.id,C.company_name,C.ticker_url FROM companies AS C INNER JOIN company_sector AS CS ON CS.company_id=C.id WHERE CS.sector_id=?"
	companies := make(map[uint64]Company)
	raws, err := mySQL.db.Query(selectQuery, id)
	if err != nil {
		return nil, err
	}
	defer raws.Close()
	for raws.Next() {
		var company Company
		var tickerUrl sql.NullString
		err = raws.Scan(&company.ID, &company.Name, &tickerUrl)
		if err != nil {
			configs.Logger.Errorf("error in scaning companies:%s", err.Error())
			continue
		}
		company.TickerUrl = tickerUrl.String
		companies[company.ID] = company
	}
	return companies, nil

}
func (mySQL *MySQL) InsertExchangeFilingDoc(companyID int, publishedIn, docURL string) error {
	insertQuery := "INSERT into exchange_filing_docs(company_id,published_in,doc_url)VALUES(?,?,?)"
	_, err := mySQL.db.Exec(insertQuery, companyID, publishedIn, docURL)
	return err
}
func (mySQL *MySQL) GetSectorExchangeFiling(sectorID int) (map[int64][]ExchangeFilingDoc, error) {
	selectQuery := "SELECT E.id,E.company_id,E.published_in,E.doc_url FROM exchange_filing_docs as E INNER JOIN company_sector as CS ON E.company_id=CS.company_id WHERE CS.sector_id=?"
	exchangeFilingDocMap := make(map[int64][]ExchangeFilingDoc)
	raws, err := mySQL.db.Query(selectQuery, sectorID)
	if err != nil {
		return nil, err
	}
	defer raws.Close()
	for raws.Next() {
		var exchangeFiling ExchangeFilingDoc
		var docURL sql.NullString
		err = raws.Scan(&exchangeFiling.ID, &exchangeFiling.CompanyID, &exchangeFiling.PublishedIn, &docURL)
		if err != nil {
			configs.Logger.Errorf("error in scaning exchange-filings:%s", err.Error())
			continue
		}
		exchangeFiling.DocURL = docURL.String
		_, isFound := exchangeFilingDocMap[exchangeFiling.CompanyID]
		if !isFound {
			exchangeFilingDocMap[exchangeFiling.CompanyID] = []ExchangeFilingDoc{}
		}
		exchangeFilingDocMap[exchangeFiling.CompanyID] = append(exchangeFilingDocMap[exchangeFiling.CompanyID], exchangeFiling)
	}
	return exchangeFilingDocMap, nil

}

func (mySQL *MySQL) InsertCompany(sectorIDs []int, company, tickerURL string) error {
	insertQuery := "INSERT INTO companies(company_name,ticker_url) VALUES(?,?)"
	transaction, err := mySQL.db.Begin()
	if err != nil {
		transaction.Rollback()
		return err
	}
	result, err := transaction.Exec(insertQuery, company, tickerURL)
	if err != nil {
		transaction.Rollback()
		return err
	}
	companyID, err := result.LastInsertId()
	if err != nil {
		transaction.Rollback()
		return err
	}
	err = insertCompanySector(transaction, sectorIDs, companyID)
	if err != nil {
		transaction.Rollback()
		return err
	}
	return transaction.Commit()
}

func insertCompanySector(transaction *sql.Tx, sectorIDs []int, companyID int64) error {
	insertQuery := "INSERT INTO company_sector(sector_id,company_id)VALUES(?,?)"
	for _, sectorID := range sectorIDs {
		_, err := transaction.Exec(insertQuery, sectorID, companyID)
		if err != nil {
			transaction.Rollback()
			return err
		}
	}
	return nil
}
