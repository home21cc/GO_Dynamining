package controllers

import (
	"dynamining/setting"
	"dynamining/models"
	"dynamining/dtools"
	"net/http"
	"github.com/gorilla/securecookie"
	"time"
	"fmt"
)

var (
	tempRoot string				// template Root
	basePage string				// base page
	basicPage string			// basic page
	runPage string				// run page
	tSysUser *models.TSysUser	// User Information
	dynaCookie = securecookie.New([]byte(dtools.GenerateUniqueId(64)), []byte(dtools.GenerateUniqueId(32)))
)

func init() {
	tempRoot = setting.TemplateConfig.TemplatesRoot
	basePage = setting.TemplateConfig.BasePage
	basicPage = setting.TemplateConfig.BasicPage
}

func setSession(idToken string, w http.ResponseWriter) {
	value := map[string]string{
		"IdToken": idToken,
	}
	if encoded, err := dynaCookie.Encode("Session", value); err == nil {
		cookie := &http.Cookie{
			Name: "Session",
			Value: encoded,
			Path: "/",
			/*Setting of Expires time
			Expires - time type 으로 선언
			Default - Session base : 창이 닫길때 (현재 Default 기능 사용중  )*/
			Expires: time.Now().Add(time.Duration(setting.AppConfig.SessionExpireTime)*time.Minute),
			/* MaxAge int type 으로 선언 ( 현재시간 + int(Sec)
			MaxAge: int(setting.AppConfig.SessionExpireTime),
			MaxAge: int(setting.AppConfig.SessionExpireTime*60),	// 분단위*/
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}
}

func getSession(r *http.Request) (idToken string) {
	if cookie, err := r.Cookie("Session"); err == nil {
		cookieValue := make(map[string]string)
		if err = dynaCookie.Decode("Session", cookie.Value, &cookieValue); err == nil {
			idToken = cookieValue["IdToken"]
		}
	}
	fmt.Println(idToken)
	return idToken
}

func renewalSession(idToken string, w http.ResponseWriter, r *http.Request) {
	idToken = getSession(r)
	setSession(idToken, w)
}

func clearSession(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name: "Session",
		Value: "",
		Path: "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}