package main

import (
	"github.com/tickrlytics/tickerlytics-backend/configs"
	"github.com/tickrlytics/tickerlytics-backend/handlers"
)

func main() {
	handlers.SetUpRoutes()
	configs.Logger.Info("Tickerlytics service start")
}
