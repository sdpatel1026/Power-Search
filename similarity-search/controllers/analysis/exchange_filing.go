package analysis

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tickrlytics/tickerlytics-backend/configs"
	"github.com/tickrlytics/tickerlytics-backend/controllers"
	"github.com/tickrlytics/tickerlytics-backend/models"
	"github.com/tickrlytics/tickerlytics-backend/services"
)

func UploadExchangeFiling(gContext *gin.Context) {
	defer controllers.Recover(gContext)
	requestID := uuid.New().String()
	configs.Logger.Infof("req-id : %s : Got an request for uploading exchange filing", requestID)
	companyIDStr := gContext.PostForm(configs.KEY_COMPANY_ID)
	companyID, err := strconv.Atoi(companyIDStr)
	if err != nil {
		controllers.ProcessError(gContext, requestID, configs.INVALID_COMPANY_ID, http.StatusBadRequest, configs.STATUS_ERROR)
		return
	}
	docFileHeader, err := gContext.FormFile(configs.KEY_DOC)
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		controllers.ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	publishedIn := gContext.PostForm(configs.KEY_PUBLISHED_IN)
	isMatched, err := regexp.MatchString(configs.PUBLISHED_IN_REGEX, publishedIn)
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		controllers.ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	} else if !isMatched {
		controllers.ProcessError(gContext, requestID, configs.INVALID_PUBLISHED_IN, http.StatusBadRequest, configs.STATUS_BAD_REQUEST)
		return
	}
	azureBlob, err := services.GetAzureBlob()
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		controllers.ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	// ioutil.ReadAll(docFileHeader.)
	doc, err := docFileHeader.Open()
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		controllers.ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	docContent, err := io.ReadAll(doc)
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		controllers.ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	containerName := configs.GetEnvWithKey(configs.KEY_EXCHANGE_FILING_DOCS_BLOB_CONTAINER, "")
	fileName := fmt.Sprintf("%d_%s", time.Now().UTC().Unix(), docFileHeader.Filename)
	docUrl, err := azureBlob.UploadDoc(docContent, controllers.TYPE_PDF, fileName, containerName)
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		//needs to classify error and send error message and status accordingly
		controllers.ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	err = models.GetMySql().InsertExchangeFilingDoc(companyID, publishedIn, docUrl)
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		//needs to classify error and send error message and status accordingly
		controllers.ProcessError(gContext, requestID, err.Error(), http.StatusInternalServerError, configs.STATUS_ERROR)
		azureBlob.DeleteDoc(fileName, containerName)
		return
	}
	controllers.Respond(gContext, requestID, configs.DOCUMENT_UPLOADED, http.StatusOK, configs.STATUS_SUCCESS)

}

func SectorExchangeFilings(gContext *gin.Context) {
	defer controllers.Recover(gContext)
	requestID := uuid.New().String()
	configs.Logger.Infof("req-id : %s : Got an request for retreiving exchange filing docs", requestID)
	sectorIDStr := gContext.Param(configs.KEY_SECTOR_ID)
	sectorID, err := strconv.Atoi(sectorIDStr)
	if err != nil {
		configs.Logger.Errorf("error in converting sector_id to integer:%s", err.Error())
		controllers.ProcessError(gContext, requestID, configs.INVALID_SECTOR_ID, http.StatusBadRequest, configs.STATUS_BAD_REQUEST)
		return
	}
	exchangeFilings, err := models.GetMySql().GetSectorExchangeFiling(sectorID)
	if err != nil {
		configs.Logger.Errorf("req-id : %s : error in retrieving exchange-filings for the sector-id %d is %s", requestID, sectorID, err.Error())
		controllers.ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	controllers.Respond(gContext, requestID, exchangeFilings, http.StatusOK, configs.STATUS_SUCCESS)

}
