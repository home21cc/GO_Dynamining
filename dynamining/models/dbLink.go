package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"dynamining/setting"
	"dynamining/dtools"
)

var MySqldb *sql.DB

func init() {
	var err error
	db, err := sql.Open("mysql", setting.AppConfig.MySqlDBUser + ":" + setting.AppConfig.MySqlDBPwd + "@tcp(" +
		setting.AppConfig.MySqlDBHost + ":" + setting.AppConfig.MySqlDBPort + ")/" + setting.AppConfig.MySqlDatabase)
	if (err != nil ) {
		dtools.Info("mdoel.dbLink.init:", err)
	}
	err = db.Ping()
	if err != nil {
		dtools.Info("model.dbLink.init:", err)
	}
	MySqldb = db
}
