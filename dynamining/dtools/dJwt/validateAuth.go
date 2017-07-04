package dJwt


import (
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"fmt"
)

func ValidateTokenAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	authBackend := InitJWTAuth()

	//token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
	// tokenstring, claims, keyFunc
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		} else {
			return authBackend.PublicKey, nil
		}
	})
	fmt.Println("Validate token:", token)
	if err == nil && token.Valid && !authBackend.IsInBlackList(r.Header.Get("Authorization")) {
		next(w, r)
	} else {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
	}
}
