package models

import (
	"time"
)

type TSysUser struct {
	Id 				string			`json:"Id"`						// 사용자 아이디
	Password		string			`json:"Password"`				// 비밀번호
	HashPass 		string			`json:"HashPass"`				// HASHED PASSWORD
	IPAddress		string			`json:"IPAddress"`		 		// 접속자 IP address
	Name 			string			`json:"Name"`					// 사용자 이름
	Company 		string			`json:"Company"`				// 사용자 회사
	Enable 			string			`json:"Enable"`					// 사용자 가능
	ReturnCode		error			`json:"ReturnCode"`				// 오류 코드 리턴
	ReturnValue		string			`json:"RuturnValue"`			// 오류 값 리턴
	CDate 			time.Time		`json:"CreateDate"`				// 생성일
	UDate 			time.Time		`json:"UpdateDate"`				// 수정일
	CUser 			string			`json:"CreateUser"`				// 생성자
	UUser 			string			`json:"UpdateUser"`				// 수정자
}
