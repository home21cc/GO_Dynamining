package controllers

import (
	"net/http"
	"dynamining/setting"
	"dynamining/models/services"
	"dynamining/models"
)

// 시작 페이지
func Start(w http.ResponseWriter, r *http.Request) {
	sessionIdToken := getSession(r)
	switch r.Method {
	case "GET":
		if sessionIdToken == "" {
			startRendering(w, r)
	 	} else {
			renewalSession(sessionIdToken, w, r)
			runPage = setting.TemplateConfig.InformationPage
			rendering(w, r, runPage, basePage, nil)
		}
	case "POST":
		tSysUser = new(models.TSysUser)
		tSysUser.Id = r.PostFormValue("inputEmail")
		tSysUser.Password = r.PostFormValue("inputPassword")
		tSysUser.IPAddress = r.Host

		// 1) id, password client to server and receive idToken server to client
		status, idToken := services.AuthenticationUser(tSysUser)   // User authentication

		// 2) save idtoken in to Session
		setSession(string(idToken), w)

		if status == http.StatusOK && idToken != nil {
			runPage = setting.TemplateConfig.InformationPage
			rendering(w, r, runPage, basePage, nil)
		} else {
			startRendering(w, r)
		}
	}

}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	/*
	dynaSession, err := dynaSessionsStore.Get(r, dynaSessionName)
	if err != nil {
		startRendering(w)

	} else {
		//userId := session.Values["UserId"]

		fmt.Println(dynaSession.Values["idToken"])
		//session.Values["idToken"] = string(services.RefreshIdToken())
		dynaSession.Save(r, w)
	}
	*/
}

// The Logout function needs to be checked
func Logout(w http.ResponseWriter, r *http.Request) {
	clearSession(w)
	startRendering(w, r)
}
