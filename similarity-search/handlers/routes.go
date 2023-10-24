package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/tickrlytics/tickerlytics-backend/configs"
	"github.com/tickrlytics/tickerlytics-backend/controllers"
	"github.com/tickrlytics/tickerlytics-backend/controllers/analysis"
)

var router *gin.Engine

func SetUpRoutes() {
	router = gin.New()
	release1 := router.Group("/r1", CORS)

	release1.GET("/discovery/docs", controllers.DocSearch)
	release1.GET("/sectors", controllers.Sectors)
	release1.GET("/companies", controllers.Companies)
	release1.GET("/sector_companies/:sector_id", controllers.SectorCompanies)
	release1.GET("/analysis/exchange_filings/:sector_id", analysis.SectorExchangeFilings)

	release1.POST("/company", AuthenticateAdmin, controllers.UploadCompany)
	release1.POST("/discovery/doc", AuthenticateAdmin, controllers.UploadDiscoveryDoc)
	release1.POST("/discovery/train", AuthenticateAdmin, controllers.TrainDiscoveryDocs)
	release1.POST("/analysis/exchange_filing", AuthenticateAdmin, analysis.UploadExchangeFiling)

	serverAddress := fmt.Sprintf(":%s", configs.GetEnvWithKey("APP_PORT", ""))
	err := router.Run(serverAddress)
	if err != nil {
		configs.Logger.Fatalf("error in start server: %s", err.Error())
	}
}
