package dtools

import (
	"database/sql"
	"strconv"
	"strings"

)

// Rows 의 Count
func RowsCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			//logError(err)
			Info(err)
		}
	}
	return count
}

func CovertBytesToString(b []byte) string {
	s := make([]string, len(b))
	for i := range b {
 		s[i] = strconv.Itoa(int(b[i]))
	}
	return strings.Join(s, " ")
}



/*
type ColName struct {
	CIndex int				// 컬럼 순서
	CName string			// 컬럼 이름
}

// Column name Return
func RowsCName(rows *sql.Rows) (){
	// colNames : 컬럼 집합
	colNames, err := rows.Columns()
	if err != nil {
		logError(err)
	}


	colNames = make([]ColName, 0, 1)

	for rows.Next() {
		colName := ColName{}
		pointers := make([]interface{}, len(colNames))
		structValue := reflect.ValueOf(colName)
		for i, colName := range colNames {
			fieldVal := structValue.FieldByName(strings.Title(colName))

			if !fieldVal.IsValid() {
				log.Fatal("field not valid")
			}
			pointers[i] = fieldVal.Addr().Interface()
		}
		err := rows.Scan(pointers...)
		if err != nil {
			logError(err)
		}
		colNames = append(colNames, colName)
	}

	return colNames
}


func logError(err error) {
	if err != nil {
		panic(err)
	}
}
*/