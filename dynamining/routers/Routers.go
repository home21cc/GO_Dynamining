package routers

import (
	"github.com/gorilla/mux"
	"net/http"
	"dynamining/dtools"
)

// 함수 호출을 위한 선언
type HandleFunc func(http.ResponseWriter, *http.Request)

// Router initialize
func InitRoutes() *mux.Router {
	dtools.Info("[====================================][InitRouters][====================================]")
	router := mux.NewRouter().StrictSlash(true)
	router = SetStartRouters(router)					// Start Page
	router = SetInformationRouters(router)				// Index Page
	return router
}

// System Panic 이 발생했을때 안정적으로 운영하기 위하여 LOG 를 남기고 운영
func robustPanics(function HandleFunc) HandleFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 접속자 추적
		dtools.Info("Connection RemoteAddress : ", r.RemoteAddr)
		dtools.Info("Request Url : ", r.RequestURI, r.Method, r.URL)

		// Panic 발생시 Recover &
		defer func() {
			err := recover()
			if err != nil {
				dtools.Critical("\n")
				dtools.Critical("[====================================][RobustPanics][====================================]")
				dtools.Critical("[", r.RemoteAddr, "]", err)
			}
		}()
		function(w, r)
	}
}
