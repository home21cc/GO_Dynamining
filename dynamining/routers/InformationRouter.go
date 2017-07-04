package routers

import (
	"github.com/gorilla/mux"
	"dynamining/controllers"
)

func SetInformationRouters(router *mux.Router) *mux.Router {
	// Information
	router.HandleFunc("/Information", robustPanics(controllers.Information)).Methods("GET", "SET")
	return router
}
