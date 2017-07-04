package logs

import (
	"runtime"
	"encoding/json"
	"os"
	"time"
)

type Brush func(string) string

func NewBrush(color string) Brush{
	pre 	:= "\033["
	reset 	:= "\033[0m"
	return func(text string) string {
		return pre + color + "m" + text + reset
	}
}

var colors = []Brush{
	NewBrush("1;37"),		// Emergency 		--> white
	NewBrush("1;36"),		// Alert 			--> cyan
	NewBrush("1;35"),		// Critical 		--> magenta
	NewBrush("1;31"),		// Error 			--> red
	NewBrush("1;33"),		// Warning 			--> yellow
	NewBrush("1;32"),		// Notice 			--> green
	NewBrush("1;34"),		// Informational 	--> blue
	NewBrush("1;34"),		// Debug 			--> blue
}

type LogConsoleWriter struct {
	logConsoleAdapter 	*logWriter
	Level 				int `json:"Loglevel"`
	Colorful			bool `json:"Color"`
}


func newConsoleWriter() LogInterface {
	logConsole := &LogConsoleWriter{
		logConsoleAdapter:			newLogWriter(os.Stdout),
		Level: 		LevelDebug,
		Colorful:	runtime.GOOS != "windows",
	}
	return logConsole
}

func (logConsole *LogConsoleWriter) Init(jsonConfig string) error {
	if len(jsonConfig) == 0 {
		return nil
	}
	err :=json.Unmarshal([]byte(jsonConfig), logConsole)
	if runtime.GOOS == "windows" {
		logConsole.Colorful = false
	}
	return err
}

func (logConsole *LogConsoleWriter) WriteMessage(when time.Time, logMessage string, level int) error {
	if level > logConsole.Level{
		return nil
	}
	if logConsole.Colorful {
		logMessage = colors[level](logMessage)
	}
	logConsole.logConsoleAdapter.println(when, logMessage)
	return nil
}

func (logConsole *LogConsoleWriter) Destroy() {

}

func (logConsole *LogConsoleWriter) Flush(){

}

func init() {
	Register(AdapterConsole, newConsoleWriter)
}