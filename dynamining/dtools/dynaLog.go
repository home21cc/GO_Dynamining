
package dtools

import (
	"strings"
	"dynamining/dtools/logs"
	"fmt"
)

const (
	LevelEmergency = iota
	LevelAlert
	LevelCritical
	LevelError
	LevelWarning
	LevelNotice
	LevelInformational
	LevelDebug
)

var DynaLog = logs.GetDynamicLog()


func SetLevel(level int) {
	logs.SetLevel(level)
}

func SetLogger(adapterName string, jsonConfig string) error {
	return logs.SetLogger(adapterName, jsonConfig)
}

func SetLogFuncCall(b bool) {
	logs.SetLogFuncCall(b)
}


func Emergency(v ...interface{}) {
	fmt.Println(v...)
	logs.Emergency(generateFmtString(len(v)), v...)
}

func Alert(v ...interface{}) {
	fmt.Println(v...)
	logs.Alert(generateFmtString(len(v)), v...)
}

func Critical(v ...interface{}) {
	fmt.Println(v...)
	logs.Critical(generateFmtString(len(v)), v...)
}

func Error(v ...interface{}) {
	fmt.Println(v...)
	logs.Error(generateFmtString(len(v)), v...)
}

func Warning(v ...interface{}) {
	fmt.Println(v...)
	logs.Warning(generateFmtString(len(v)), v...)
}

func Warn(v ...interface{}) {
	fmt.Println(v...)
	logs.Warn(generateFmtString(len(v)), v...)
}

func Notice(v ...interface{}) {
	fmt.Println(v...)
	logs.Notice(generateFmtString(len(v)), v...)
}

func Informational(v ...interface{}) {
	fmt.Println(v...)
	logs.Informational(generateFmtString(len(v)), v...)
}

func Info(v ...interface{}) {
	fmt.Println(v...)
	logs.Info(generateFmtString(len(v)), v...)
}

func Debug(v ...interface{}) {
	fmt.Println(v...)
	logs.Debug(generateFmtString(len(v)), v...)
}

func Trace(v ...interface{}) {
	fmt.Println(v...)
	logs.Trace(generateFmtString(len(v)), v...)
}

func generateFmtString(n int) string {
	return strings.Repeat("%v ", n)
}

