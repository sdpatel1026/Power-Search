package controllers

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tickrlytics/tickerlytics-backend/configs"
	"github.com/tickrlytics/tickerlytics-backend/models"
	"github.com/tickrlytics/tickerlytics-backend/services"
)

func Companies(gContext *gin.Context) {
	defer Recover(gContext)
	requestID := uuid.New().String()
	configs.Logger.Infof("req-id : %s : Got an request for retrieving companies", requestID)
	companies, err := models.GetMySql().GetCompanies()
	if err != nil {
		configs.Logger.Infof("req-id: %s :error in retrieving companies %s", requestID, err.Error())
		ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	Respond(gContext, requestID, companies, http.StatusOK, configs.STATUS_SUCCESS)
}

func UploadCompany(gContext *gin.Context) {
	defer Recover(gContext)
	requestID := uuid.New().String()
	configs.Logger.Infof("req-id : %s : Got an request for uploading companies", requestID)
	companyName := gContext.PostForm(configs.KEY_COMPANY)
	if companyName == "" {
		ProcessError(gContext, requestID, configs.COMPANY_REQUIRED, http.StatusBadRequest, configs.STATUS_BAD_REQUEST)
		return
	}
	sectorIDs := gContext.PostFormArray(configs.KEY_SECTOR_IDS)
	if len(sectorIDs) == 0 {
		ProcessError(gContext, requestID, configs.SECTOR_IDS_REQUIRED, http.StatusBadRequest, configs.STATUS_BAD_REQUEST)
		return
	}
	var sectorIDsInt []int
	for _, sectorID := range sectorIDs {
		sectorIDInt, err := strconv.Atoi(sectorID)
		if err != nil {
			ProcessError(gContext, requestID, configs.SECTOR_IDS_MUST_BE_INT, http.StatusBadRequest, configs.STATUS_BAD_REQUEST)
			return
		}
		sectorIDsInt = append(sectorIDsInt, sectorIDInt)
	}
	tickerFileHeader, err := gContext.FormFile(configs.KEY_TICKER)
	if err != nil {
		configs.Logger.Errorf("error in parsing ticker file: %s", err.Error())
		ProcessError(gContext, requestID, configs.TICKER_REQUIRED, http.StatusBadRequest, configs.STATUS_BAD_REQUEST)
		return
	}

	azureBlob, err := services.GetAzureBlob()
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}

	ticker, err := tickerFileHeader.Open()
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	tickerContent, err := io.ReadAll(ticker)
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	if !isImage(tickerContent) {
		configs.Logger.Errorf("req-id : %s \n:file is not a pdf\n", requestID)
		ProcessError(gContext, requestID, configs.INVALID_TICKER_FORMAT, http.StatusBadRequest, configs.STATUS_BAD_REQUEST)
	}
	containerName := configs.GetEnvWithKey(configs.KEY_TICKER_BLOB_CONTAINER, "")
	fileName := fmt.Sprintf("%d_%s", time.Now().UTC().Unix(), tickerFileHeader.Filename)
	tickerURL, err := azureBlob.UploadDoc(tickerContent, getContentType(tickerContent), fileName, containerName)
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		//needs to classify error and send error message and status accordingly
		ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	err = models.GetMySql().InsertCompany(sectorIDsInt, companyName, tickerURL)
	if err != nil {
		configs.Logger.Errorf("req-id : %s : %s", requestID, err.Error())
		azureBlob.DeleteDoc(fileName, containerName)
		//needs to classify error and send error message and status accordingly
		ProcessError(gContext, requestID, err.Error(), http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	Respond(gContext, requestID, configs.COMPANY_SUCCESSFULLY_UPLOADED, http.StatusOK, configs.STATUS_SUCCESS)
}
