package main

import (
	"net/http"
	"runtime"
	"github.com/codegangsta/negroni"
	"dynamining/dtools"
	"dynamining/routers"
	"dynamining/setting"
	"dynamining/models"

)

// Program name : DynaMining
// Program Description : Data mining 을 활용한 수주 자료 분석
// Developer : DoMyoung, Park
// Create Date : 2017/04/30

func main() {
	initialize()
	router := routers.InitRoutes()
	defer models.MySqldb.Close()
	mux := negroni.Classic()
	mux.UseHandler(router)

	// Http Server Parsing
	server := &http.Server {
		Addr: setting.AppConfig.ServerIP,
		Handler: mux,

	}
	if setting.AppConfig.HttpTLS != "TLS" {
		server.ListenAndServe()
	}  else {
		// 사용중이지 않음
		//server.ListenAndServeTLS("cert.pem", "key.pem")
	}
}

func initialize() {

	if runtime.NumCPU() == 1 {
		runtime.GOMAXPROCS(1)
	}  else {
		runtime.GOMAXPROCS(runtime.NumCPU()-1)			// 실행 중인 시스템의 CPU 갯수 - 1을 사용하도록
	}

	dtools.SetLevel(dtools.LevelInformational)		// Log Level 6레벨로 지정
	dtools.DynaLog.SetDynamicLog(setting.AppConfig.LogAdapter, `{"Filename":"Configs/logs.log"}` )
	dtools.Info("\r\n")
	dtools.Info("[======================================== START ========================================]")
	dtools.Info("  - CPU Number : ", (runtime.NumCPU()) )
	dtools.Info("  - Running CPU Number : ", (runtime.NumCPU()-1) )
	dtools.Info("  - Version : ", setting.AppConfig.Version)
	dtools.Info("  - Server IP : ", setting.AppConfig.ServerIP)
	dtools.Info("  - LogAdapter : ", setting.AppConfig.LogAdapter)
	dtools.Info("  - JWTExpireTime: ", setting.AppConfig.JWTExpireTime, "minute")
	dtools.Info("  - SessionExpireTime: ", setting.AppConfig.SessionExpireTime, "minute")
}
