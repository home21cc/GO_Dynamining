package dJwt

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"time"
	"dynamining/dtools"
	"os"
	"dynamining/setting"
	"bufio"
	"encoding/pem"
	"crypto/x509"
	"dynamining/models"
	"golang.org/x/crypto/bcrypt"
)

// Private, Public Key
type JWTAuth	struct {
	privateKey	*rsa.PrivateKey
	PublicKey	*rsa.PublicKey
}

var authInstance *JWTAuth = nil

func InitJWTAuth() *JWTAuth {
	if authInstance == nil {
		authInstance = &JWTAuth{
			privateKey: getPrivateKey(),
			PublicKey: getPublicKey(),
		}
	}
	return authInstance
}

// Token 생성
// index description
// Audience  string `json:"aud,omitempty"`  : 수신자
// ExpiresAt int64  `json:"exp,omitempty"`  : 만료시간
// Id        string `json:"jti,omitempty"`  : jwt id
// IssuedAt  int64  `json:"iat,omitempty"`  : 발급시간
// Issuer    string `json:"iss,omitempty"`  : 발급자
// NotBefore int64  `json:"nbf,omitempty"`  : 지정시간까지 처리하지 않아야 함
// Subject   string `json:"sub,omitempty"`  : Token 화 대상
func (jwtAuth *JWTAuth) GenerateToken(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	/*
	token.Claims = &DynamicToken{
		&jwt.StandardClaims{
			// 만료 시간
			ExpiresAt: time.Now().Add(time.Duration(setting.AppConfig.JWTExpireTime)*time.Second).Unix(),
			// 발급 시간
			IssuedAt: time.Now().Unix(),
			// Token화 대상
			Subject: userId,
		},
		//"Level 1",
		CustomerInfo{userId},
	}
	*/
	token.Claims = jwt.MapClaims{
		"exp": time.Now().Add(time.Duration(setting.AppConfig.JWTExpireTime)*time.Minute).Unix(),
		"iat": time.Now().Unix(),
		"sub": userId,
	}
	tokenString, err := token.SignedString(jwtAuth.privateKey)
	if err != nil {
		dtools.Critical("dJwt.jwtAuthentication.GenerateToken: Generate token error")
		panic(err)
		return "", err
	}
	return tokenString, nil
}
// Hash password
func (jwtAuth *JWTAuth) Authenticate(tSysUser *models.TSysUser) (string) {
	//hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Dynamining"), bcrypt.DefaultCost)
	password := []byte(tSysUser.Id + tSysUser.Password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(password, 10)
	return string(hashedPassword)
}

func (jwtAuth *JWTAuth) CompareHashPassword(tSysUser *models.TSysUser, hashPass string) bool {
	tPass := []byte(tSysUser.Id + tSysUser.Password)
	hPass := []byte(hashPass)
	err := bcrypt.CompareHashAndPassword(hPass, tPass)

	return err == nil
}

func (jwAuth *JWTAuth) GenerateHashPassword(tSysUser *models.TSysUser) (string) {
	// Id, Password 를 조합하여 처리
	password := []byte(tSysUser.Id + tSysUser.Password)
	hPassword, _ := bcrypt.GenerateFromPassword(password, 12)
	return string(hPassword)
}

// -------------------------------------------------------------------------------------------------
//type Keyfunc func(*DynamicToken) (interface{}, error)
/*
// Define Customer
type CustomerInfo struct {
	Id 		string
	// Kind 	string							// 용도 확인 불가
}

// Define Token
type DynamicToken struct {
	*jwt.StandardClaims
	// TokenType 	string						// 용도 확인불가
	CustomerInfo
}

*/

func  (jwtAuth *JWTAuth) getTokenRemainingValidity(timestamp interface{}) int64 {
	if validity, ok := timestamp.(float64); ok {
		tm := time.Unix(int64(validity), 0)
		remaining := tm.Sub(time.Now())
		if remaining > 0 {
			return int64(remaining.Seconds()) + int64(setting.AppConfig.ExpireOffset)
		}
	}
	return setting.AppConfig.ExpireOffset
}

func(jwtAuth *JWTAuth) Logout(tokenString string, token *jwt.Token) error {
	redisConn := Connect()
	return redisConn.SetValue(tokenString, tokenString, jwtAuth.getTokenRemainingValidity(token.Claims))
}

func(jwtAuth *JWTAuth) IsInBlackList(token string) bool {
	redisConn := Connect()
	redisToken, _ := redisConn.GetValue(token)

	if redisToken == nil {
		return false
	}
	return true
}

// Get Private Key
func getPrivateKey() *rsa.PrivateKey {
	privateKeyFile, err := os.Open(setting.AppConfig.PrivateKeyPath)
	if err != nil {
		dtools.Critical("dJwt.jwtAuth.getPrivateKey: Can't find privateKey File ", setting.AppConfig.PrivateKeyPath)
		panic(err)
	}
	permissionInfo, _ := privateKeyFile.Stat()

	size := permissionInfo.Size()
	permissionBytes := make([]byte, size)

	buffer := bufio.NewReader(privateKeyFile)
	_, err = buffer.Read(permissionBytes)
	data, _ := pem.Decode([]byte(permissionBytes))

	privateKeyFile.Close()
	privateKeyImported, err := x509.ParsePKCS1PrivateKey(data.Bytes)
	if err != nil {
		dtools.Critical("dJwt.jwtAuth.getPrivateKey: Can't make x509 Privatekey")
		panic(err)
	}
	return privateKeyImported
}

// Get Public Key
func getPublicKey() *rsa.PublicKey {
	publicKeyFile, err := os.Open(setting.AppConfig.PublicKeyPath)
	if err != nil {
		dtools.Critical("dJwt.jwtAuth.getPublicKey: Can't fild publicKey File : ", setting.AppConfig.PublicKeyPath)
		panic(err)
	}
	permissionInfo, _ := publicKeyFile.Stat()
	size := permissionInfo.Size()
	permissionBytes := make([]byte, size)
	buffer := bufio.NewReader(publicKeyFile)
	_, err = buffer.Read(permissionBytes)
	data, _ := pem.Decode([]byte(permissionBytes))

	publicKeyFile.Close()
	publicKeyImported, err := x509.ParsePKIXPublicKey(data.Bytes)
	if err != nil {
		dtools.Critical("dJwt.jwtAuth.getPublicKey: Can't make x509 Publickey")
		panic(err)
	}
	rsaPub, ok := publicKeyImported.(*rsa.PublicKey)

	if !ok {
		panic(err)
	}

	return rsaPub
}