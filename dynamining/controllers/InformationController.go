package controllers

import (
	"net/http"
	"dynamining/setting"

	"fmt"
)

func Information(w http.ResponseWriter, r *http.Request) {
	runPage = setting.TemplateConfig.InformationPage
	w.Header().Add("Location", "/Information")

	idToken := getSession(r)
	if idToken == "" {
		startRendering(w, r)
		fmt.Println("r.Header", r.Header)
	} else {
		rendering(w, r, runPage, basePage, nil)
		fmt.Println("r.Header", r.Header)
	}
}
