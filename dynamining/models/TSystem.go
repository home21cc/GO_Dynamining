package models

import "time"

type (
	TCode struct {
		CodeMajor 		string			`json:"CodeMajor"`
		CodeMinor		string			`json:"CodeMinor"`
		FirstValue		string			`json:"FirstValue"`
		SecondValue		string  		`json:"SecondValue"`
		Description 	string  		`json:"Description"`
		Enable			string			`json:"Enable"`
		CDate			time.Time		`json:"CreateDate"`
		UDate 			time.Time		`json:"UpdateDate"`
		CUser 			string			`json:"CreateUser"`
		UUser 			string			`json:"UpdateUser"`
	}

)
