package controllers

import (
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tickrlytics/tickerlytics-backend/configs"
	"github.com/tickrlytics/tickerlytics-backend/models"
)

type Response struct {
	Status    int         `json:"status"`
	Result    interface{} `json:"result"`
	Error     interface{} `json:"error"`
	RequestID string      `json:"req_id"`
}
type TrainDocResponse struct {
	DocID   uint64 `json:"doc_id"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func init() {
	err := models.Connect()
	if err != nil {
		configs.Logger.Fatalf("error in connecting to DB:%s", err.Error())
	}
}

func ProcessError(gContext *gin.Context, reqID string, err interface{}, httpStatusCode, statusCode int) {
	response := Response{}
	response.Status = statusCode
	response.Error = err
	response.RequestID = reqID
	gContext.AbortWithStatusJSON(httpStatusCode, response)
}
func Respond(gContext *gin.Context, reqID string, result interface{}, httpStatusCode, statusCode int) {
	response := Response{}
	response.Status = statusCode
	response.Result = result
	response.RequestID = reqID
	gContext.JSON(httpStatusCode, response)
}
func Recover(ctx *gin.Context) {
	if err := recover(); err != nil {
		configs.Logger.Errorf("exit crawl: PANIC occured - %v", err)
		replacer := strings.NewReplacer("\n", " >> ", "\t", " ")
		configs.Logger.Errorf("stacktrace from panic: %s", replacer.Replace(string(debug.Stack())))
		ctx.JSON(http.StatusInternalServerError, Response{Status: configs.STATUS_ERROR, Error: configs.TECHNICAL_ERROR})

	}
}
