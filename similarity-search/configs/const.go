package configs

const (
	KEY_LOGFILE_PATH                        string = "LOGFILE_PATH"
	KEY_LOGFILE_ENCODING                           = "LOGFILE_ENCODING"
	KEY_AZURE_STORAGE_ACCOUNT_KEY                  = "AZURE_ACCOUNT_KEY"
	KEY_AZURE_STORAGE_ACCOUNT_NAME                 = "AZURE_ACCOUNT_NAME"
	KEY_DISCOVERY_DOCS_BLOB_CONTAINER              = "DISCOVERY_DOCS_BLOB_CONTAINER"
	KEY_EXCHANGE_FILING_DOCS_BLOB_CONTAINER        = "EXCHANGE_FILING_DOCS_BLOB_CONTAINER"
	KEY_TICKER_BLOB_CONTAINER                      = "TICKER_BLOB_CONTAINER"
	KEY_STATUS                                     = "status"
	KEY_ERROR                                      = "error"
	KEY_ADMIN                                      = "admin"
	KEY_COMPANY                                    = "company"
	KEY_COMPANY_ID                                 = "company_id"
	KEY_DOC_TYPE                                   = "doc_type"
	KEY_DOC                                        = "doc"
	KEY_PUBLISHED_DATE                             = "published_date"
	KEY_PUBLISHED_IN                               = "published_in"
	KEY_RESULT                                     = "result"
	KEY_RANK                                       = "rank"
	KEY_REQ_ID                                     = "req_id"
	KEY_DOCS                                       = "docs"
	KEY_DOC_ID                                     = "doc_id"
	KEY_SECTOR_ID                                  = "sector_id"
	KEY_SECTOR_IDS                                 = "sector_ids"
	KEY_DOC_NAME                                   = "doc_name"
	KEY_TEXT                                       = "text"
	KEY_MSG                                        = "msg"
	KEY_LIMIT                                      = "limit"
	KEY_OFFSET                                     = "offset"
	KEY_TICKER                                     = "ticker"

	INVALID_COMPANY_ID            = "Invalid company_id."
	INVALID_LIMIT                 = "limit must be a non-negative integer."
	INVALID_OFFSET                = "offset must be a non-negative integer."
	INVALID_SECTOR_ID             = "Invalid sector_id."
	INVALID_PUBLISHED_IN          = "Invalid published_in."
	INVALID_TICKER_FORMAT         = "Ticker must be in jpeg/png format."
	DOC_TYPE_REQUIRED             = "doc_type is required..."
	TEXT_REQUIRED                 = "text required..."
	FILE_REQUIRED                 = "docs required..."
	COMPANY_REQUIRED              = "company required..."
	SECTOR_IDS_REQUIRED           = "sector_ids required..."
	TICKER_REQUIRED               = "ticker required..."
	SECTOR_IDS_MUST_BE_INT        = "All sector_ids must be in integer."
	TECHNICAL_ERROR               = "Something went wrong. Please try again later..."
	READING_ERROR                 = "error in reading the document."
	DOCUMENT_UPLOADED             = "Document successfully uploaded."
	DOCS_TRAINED                  = "Document successfully trained."
	DOCS_ALEREADY_TRAINED         = "Documents are already trained."
	TRAINING_NOT_REQUIRED         = "This document does not needs to be trained, as it does not contains any useful information."
	ALL_DOCS_TRAINED              = "All documents are trained."
	FILE_IS_NOT_PDF               = "file is not a pdf"
	COMPANY_SUCCESSFULLY_UPLOADED = "Company successfully uploaded."
)

const (
	ADMIN                            int     = 737
	THRESHOLD_SIMILARITIES           float64 = 0.40
	STATUS_DOC_NOT_TRAINED           uint8   = 0
	STATUS_DOC_TRAINED               uint8   = 1
	STATUS_DOC_ALEREADY_TRAINED      uint8   = 2
	STATUS_DOC_TRAINING_NOT_REQUIRED uint8   = 3
	STATUS_INVALID_DOC               uint8   = 4
	STATUS_ERROR_IN_OCR              uint8   = 5
)
const (
	STATUS_ERROR int = iota
	STATUS_SUCCESS
	STATUS_BAD_REQUEST
)

// needs to remove
// const (
//
//	OUTPUT_LEN int = 10 // no. of documents given in output
//
// )
const (
	PUBLISHED_IN_REGEX = `^Q[1-4]FY[2-9][0-9]{3}$`
)
