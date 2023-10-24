package controllers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tickrlytics/tickerlytics-backend/configs"
	"github.com/tickrlytics/tickerlytics-backend/controllers/tfidf"
	"github.com/tickrlytics/tickerlytics-backend/models"
	"github.com/tickrlytics/tickerlytics-backend/services"
)

func UploadDiscoveryDoc(gContext *gin.Context) {
	defer Recover(gContext)
	requestID := uuid.New().String()
	configs.Logger.Infof("req-id : %s : Got an request for uploading doc", requestID)
	companyIDStr := gContext.PostForm(configs.KEY_COMPANY_ID)
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		ProcessError(gContext, requestID, configs.INVALID_COMPANY_ID, http.StatusBadRequest, configs.STATUS_ERROR)
		return
	}
	docType := gContext.PostForm(configs.KEY_DOC_TYPE)
	if docType == "" {
		ProcessError(gContext, requestID, configs.DOC_TYPE_REQUIRED, http.StatusBadRequest, configs.STATUS_ERROR)
		return
	}
	docFileHeader, err := gContext.FormFile(configs.KEY_DOC)
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	publishedDate := gContext.PostForm(configs.KEY_PUBLISHED_DATE)

	azureBlob, err := services.GetAzureBlob()
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	// ioutil.ReadAll(docFileHeader.)
	doc, err := docFileHeader.Open()
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	docContent, err := io.ReadAll(doc)
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	containerName := configs.GetEnvWithKey(configs.KEY_DISCOVERY_DOCS_BLOB_CONTAINER, "")
	fileName := fmt.Sprintf("%d_%s", time.Now().UTC().Unix(), docFileHeader.Filename)
	docUrl, err := azureBlob.UploadDoc(docContent, TYPE_PDF, fileName, containerName)
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		//needs to classify error and send error message and status accordingly
		ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	err = models.GetMySql().InsertDoc(companyID, docType, docUrl, publishedDate)
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		//needs to classify error and send error message and status accordingly
		ProcessError(gContext, requestID, err.Error(), http.StatusInternalServerError, configs.STATUS_ERROR)
		azureBlob.DeleteDoc(fileName, containerName)
		return
	}
	Respond(gContext, requestID, configs.DOCUMENT_UPLOADED, http.StatusOK, configs.STATUS_SUCCESS)
}

func TrainDiscoveryDocs(gContext *gin.Context) {
	defer Recover(gContext)
	requestID := uuid.New().String()
	configs.Logger.Infof("req-id : %s : Got an request to train the docs", requestID)
	offset := gContext.PostForm(configs.KEY_OFFSET)
	limit := gContext.PostForm(configs.KEY_LIMIT)
	offSetInt, err := strconv.Atoi(offset)
	if err != nil || offSetInt < 0 {
		ProcessError(gContext, requestID, configs.INVALID_OFFSET, http.StatusBadRequest, configs.STATUS_ERROR)
		return
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 0 {
		ProcessError(gContext, requestID, configs.INVALID_LIMIT, http.StatusBadRequest, configs.STATUS_ERROR)
		return
	}
	// var responseMsg string
	docs, err := models.GetMySql().GetUnTrainedDocs(requestID, offSetInt, limitInt)
	if err != nil {
		if err != sql.ErrNoRows {
			ProcessError(gContext, requestID, err.Error(), http.StatusInternalServerError, configs.STATUS_ERROR)
			return
		}
		Respond(gContext, requestID, configs.DOCS_ALEREADY_TRAINED, http.StatusOK, configs.STATUS_SUCCESS)
		return
	}
	if len(docs) == 0 {
		Respond(gContext, requestID, configs.DOCS_ALEREADY_TRAINED, http.StatusOK, configs.STATUS_SUCCESS)
		return
	}
	tfIdf, err := tfidf.New()
	if err != nil {
		ProcessError(gContext, requestID, err.Error(), http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	var trainDocsResponse []TrainDocResponse
	for _, doc := range docs {
		var trainDocResponse TrainDocResponse
		trainDocResponse.DocID = doc.ID
		configs.Logger.Infof("Trying to train doc with doc-id: %d", doc.ID)
		doc.Content, err = services.DownloadFile(doc.Url.String)
		if err != nil {
			trainDocResponse.Error = err.Error()
			trainDocsResponse = append(trainDocsResponse, trainDocResponse)
			configs.Logger.Errorf("req-id : %s \n:error in downloading doc: %s \nfor doc-id: %d", requestID, err.Error(), doc.ID)
			continue
		}
		// saveDoc(doc.Content)
		// fmt.Printf("doc.Content: %s\n", string(doc.Content))
		configs.Logger.Infof("doc file with doc-id %d downloaded successfully", doc.ID)
		if !isPDF(doc.Content) {
			trainDocResponse.Error = configs.FILE_IS_NOT_PDF
			trainDocsResponse = append(trainDocsResponse, trainDocResponse)
			configs.Logger.Errorf("req-id : %s \n:file is not a pdf\n", requestID)
			err = models.GetMySql().UpdateDocTrainStatus(doc.ID, configs.STATUS_INVALID_DOC)
			if err != nil {
				configs.Logger.Errorf("req-id : %s :error in updating train status of doc: %s\nfor doc-id: %d", requestID, err.Error(), doc.ID)
			}
			continue
		}
		content, err := services.OCR(doc.Content, int64(doc.ID))
		if err != nil {
			configs.Logger.Errorf("req-id : %s\n error in ocr %s \nfor doc-id : %d", requestID, err.Error(), doc.ID)
			trainDocResponse.Error = configs.TECHNICAL_ERROR
			trainDocsResponse = append(trainDocsResponse, trainDocResponse)
			err = models.GetMySql().UpdateDocTrainStatus(doc.ID, configs.STATUS_ERROR_IN_OCR)
			if err != nil {
				configs.Logger.Errorf("req-id : %s :error in updating train status of doc: %s\nfor doc-id: %d", requestID, err.Error(), doc.ID)
			}
			continue
		}
		content = cleanContent(content)
		doc.Content = []byte(content)
		tfIdf.TrainDoc(&doc)
		// configs.Logger.Info("doc file trained successfully")
		trainDocResponse.Message = configs.DOCS_TRAINED
		trainDocsResponse = append(trainDocsResponse, trainDocResponse)
		err = models.GetMySql().UpdateDocTrainStatus(doc.ID, doc.IsTrained)
		if err != nil {
			configs.Logger.Errorf("req-id : %s :error in updating train status of doc: %s for doc-id: %d", requestID, err.Error(), doc.ID)
		}
	}
	Respond(gContext, requestID, trainDocsResponse, http.StatusOK, configs.STATUS_SUCCESS)

}

func DocSearch(gContext *gin.Context) {
	defer Recover(gContext)
	requestID := uuid.New().String()
	configs.Logger.Infof("req-id : %s : Got an request to search the docs", requestID)
	text := strings.TrimSpace(gContext.Query(configs.KEY_TEXT))
	if text == "" {
		ProcessError(gContext, requestID, configs.TEXT_REQUIRED, http.StatusBadRequest, configs.STATUS_ERROR)
		return
	}
	tfIdf, err := tfidf.New()
	if err != nil {
		configs.Logger.Errorf("error in getting tfidf instanceL: %s", err.Error())
		ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	companyDocs := tfIdf.FindCompanyDocs(text)
	Respond(gContext, requestID, companyDocs, http.StatusOK, configs.STATUS_SUCCESS)
}
