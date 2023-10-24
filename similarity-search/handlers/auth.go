package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tickrlytics/tickerlytics-backend/configs"
)

func AuthenticateAdmin(gContext *gin.Context) {
	isAdminStr := gContext.GetHeader(configs.KEY_ADMIN)
	isAdmin, err := strconv.Atoi(isAdminStr)
	if err != nil || isAdmin != configs.ADMIN {
		gContext.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	gContext.Next()
}

func CORS(gContext *gin.Context) {
	gContext.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	gContext.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	gContext.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	gContext.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")
	gContext.Writer.Header().Set("Access-Control-Expose-Headers", "*")
	gContext.Next()
}
