package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tickrlytics/tickerlytics-backend/configs"
	"github.com/tickrlytics/tickerlytics-backend/models"
)

func Sectors(gContext *gin.Context) {
	defer Recover(gContext)
	requestID := uuid.New().String()
	configs.Logger.Infof("req-id : %s : Got an request for retrieving sectors", requestID)
	sectors, err := models.GetMySql().GetSectors()
	if err != nil {
		configs.Logger.Infof("req-id: %s :error in retrieving sectors %s", requestID, err.Error())
		ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	Respond(gContext, requestID, sectors, http.StatusOK, configs.STATUS_SUCCESS)
}

func SectorCompanies(gContext *gin.Context) {
	defer Recover(gContext)
	requestID := uuid.New().String()
	configs.Logger.Infof("req-id : %s : Got an request for retrieving sector's companies", requestID)
	sectorIDStr := gContext.Param(configs.KEY_SECTOR_ID)
	sectorID, err := strconv.Atoi(sectorIDStr)
	if err != nil {
		configs.Logger.Errorf("error in converting sector_id to integer:%s", err.Error())
		ProcessError(gContext, requestID, configs.INVALID_SECTOR_ID, http.StatusBadRequest, configs.STATUS_BAD_REQUEST)
		return
	}
	companies, err := models.GetMySql().GetSectorCompanies(int64(sectorID))
	if err != nil {
		configs.Logger.Errorf("req-id : %s : error in retrieving companies for the sector-id %d is %s", requestID, sectorID, err.Error())
		ProcessError(gContext, requestID, configs.TECHNICAL_ERROR, http.StatusInternalServerError, configs.STATUS_ERROR)
		return
	}
	Respond(gContext, requestID, companies, http.StatusOK, configs.STATUS_SUCCESS)
}
