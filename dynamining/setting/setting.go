package setting

import (
	"net/http"
	"encoding/json"
	"os"
	"log"
)


type (

	appError struct {
		Error 		string 	`json:"error"`
		Message		string 	`json:"message"`
		HttpStatus  int 	`json:"status"`
	}

	errorResource struct {
		Data appError	`json:"data"`
	}

	// System 운영을 위한 설정
	configuration struct {
		Version  string						// 프로그램 버전
		ServerIP string						// Server Internet Protocol Address

		// Database host, port, user, password, databasename
	    // Oracle Database
		OracleDBHost, OracleDBPort, OracleDBUser, OracleDBPwd, OracleDatabase string
		// MSSQL Database
		MSSQLDBHost, MSSQLDBPort, MSSQLDBUser, MSSQLDBPwd, MSSQLDatabase string
		// MySQL Database
		MySqlDBHost, MySqlDBPort, MySqlDBUser, MySqlDBPwd, MySqlDatabase string
		LogLevel int 							// Log write level
		LogAdapter string						// Log Adapter, AdapterFile, AdapterConsole, AdapterSmtp ...
		ServerLicense string					// Server License Demo, Prod ...
		PrivateKeyPath, PublicKeyPath string	// Private, Public Path
		ExpireOffset int64
		JWTExpireTime int						// Jwt token  Expiretime
		SessionExpireTime int					// Session expiretime
		HttpTLS string							// TLS Support
	}

	// Program 추가를 위한 설정
	// Program 개발시 아래 내용을 추가하고 진행할 것
	// template.json template file index 추가
	// example : "AuthenticationPage"  : "authentication",
	// templateConfig struct 에 link index 추가
	templateConfig struct {
		// URL path, file Name ...
		StaticUrl, StaticRoot string
		TemplatesUrl, TemplatesRoot string
		BasePage string						// 인증 되지 않은 페이지의 Template
		BasicPage string					// 인증 된 페이지의 Template
		StartPage string					// Start page
		AddUserPage string					// add User
		EditUserPage string					// edit User
		InformationPage string					// Index page
		Error404Page string					// Error Page
	}
)

func DisplayAppError(w http.ResponseWriter, handlerError error, message string, code int) {
	errObject := appError{
		Error:		handlerError.Error(),
		Message:	message,
		HttpStatus: code,
	}

//	Error.Printf("[AppError]: %s\n", handlerError)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if j, err := json.Marshal(errorResource{Data: errObject}); err == nil {
		w.Write(j)
	}
}

var AppConfig  configuration
var TemplateConfig templateConfig


func loadAppConfig() {
	file, err := os.Open("Configs/config.json")
	defer file.Close()

	if err != nil {
		log.Fatalf("[setting.go][loadConfig]: %s\n", err)
	}
	decoder := json.NewDecoder(file)
	AppConfig = configuration{}
	err = decoder.Decode(&AppConfig)
	if err != nil {
		log.Fatalf("[setting.go][loadAppConfig]: %s\n",err)
	}
}

func loadTemplateConfig() {
	file, err := os.Open("Configs/template.json")
	defer file.Close()

	if err != nil {
		log.Fatalf("[loadTemplateConfig]: %s\n", err)
	}
	decoder := json.NewDecoder(file)
	TemplateConfig = templateConfig{}
	err = decoder.Decode(&TemplateConfig)
	if err != nil {
		log.Fatalf("[loadTemplateConfig]: %s\n",err)
	}
}