package routers

import (
	"github.com/gorilla/mux"
	"dynamining/controllers"
	"net/http"
	"dynamining/setting"
)

func SetStartRouters(router *mux.Router) *mux.Router {
	router.HandleFunc("/", robustPanics(controllers.Start)).Methods("GET", "POST")
	//router.HandleFunc("/Refresh", robustPanics(controllers.RefreshToken)).Methods("GET")
	router.HandleFunc("/Logout", robustPanics(controllers.Logout)).Methods("GET")

	// Static file : img, css, js and etc
	// gorilla/mux Router 를 사용시에는 아래 명령을 사용해야 함
	fs := http.FileServer(http.Dir(setting.TemplateConfig.StaticRoot))
	router.PathPrefix(setting.TemplateConfig.StaticUrl).Handler(http.StripPrefix(setting.TemplateConfig.StaticUrl, fs))

	// golang Default Route 를 사용할 때는 Comment 처리된 부분 사용
	//fs := http.FileServer(http.Dir("./Static"))
	//router.Handle("/static/", http.StripPrefix("/static/", fs))
	return router
}
