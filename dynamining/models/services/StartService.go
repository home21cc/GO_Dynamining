package services

import (
	"net/http"
	"encoding/json"
	"dynamining/models"
	"dynamining/dtools"
	"dynamining/dtools/dJwt"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
)

// 사용자를 확인하고 IdToken 을 부여한다.
func AuthenticationUser(tSysUser *models.TSysUser) (int, []byte)  {
	// Init JWTAuth
	authBackend := dJwt.InitJWTAuth()
	// Database 에서 사용자가 있는지 확인
	r, err := models.MySqldb.Query("SELECT COUNT(userId) FROM tsysuser WHERE userId = '" + tSysUser.Id +
		"' AND userPass = '" + tSysUser.Password + "'")
	defer r.Close()
	if err != nil {
		dtools.Info("TSysStartService.CertificationUser: ", err)
	}
	rawValue := dtools.RowsCount(r)

	if rawValue > 0 && err == nil {
		// Id Token 생성
		idToken, err := authBackend.GenerateToken(tSysUser.Id)
		if err != nil{
			return http.StatusInternalServerError, []byte("")
		} else {
			response, _ := json.Marshal(dJwt.TokenAuth{idToken})
			return http.StatusOK, response
		}
	}  else {
		return http.StatusInternalServerError, []byte("")
	}
}
/*
func RefreshIdToken(userId string) [] byte {
	authBackend := dJwt.InitJWTAuth()
	idToken, err := authBackend.GenerateToken(userId)
	if err != nil {
		dtools.Info("TSysStartService.RefreshIdToken: ", err)
		panic(err)
	}
	response, err := json.Marshal(dJwt.TokenAuth{idToken})
	if err != nil {
		dtools.Info("TSysStartService.RefreshIdToken: ", err)
		panic(err)
	}
	return response
}
*/
func Logout(r *http.Request) error {
	authBackend := dJwt.InitJWTAuth()
	idToken, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		return authBackend.PublicKey, nil
	})

	if err != nil {
		return err
	}
	tokenString := r.Header.Get("Authorization")
	return authBackend.Logout(tokenString, idToken)
}

/*
type response struct {
	Text string `json:"Text"`
}

type token struct {
	Token string `json:"token"`
}



func jsonResponse(response interface{}, w http.ResponseWriter) {
	json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(json)
}
*/