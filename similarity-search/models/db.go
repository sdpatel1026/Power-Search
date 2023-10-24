package models

import (
	"github.com/tickrlytics/tickerlytics-backend/configs"
)

func Connect() error {
	var mySQLHost, mySQLUser, mySQLPass, mySQLDB string

	mySQLHost = configs.GetEnvWithKey("MYSQL_HOST", "")
	mySQLUser = configs.GetEnvWithKey("MYSQL_USER", "")
	mySQLPass = configs.GetEnvWithKey("MYSQL_PASS", "")
	mySQLDB = configs.GetEnvWithKey("MYSQL_DB", "")

	var err error
	mySQL = new(MySQL)
	mySQL.db, err = connectMySql(mySQLUser, mySQLPass, mySQLHost, mySQLDB)
	return err
}
