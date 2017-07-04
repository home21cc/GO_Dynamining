package logs

import (
	"encoding/json"
	"errors"
	"net/url"
	"net"
	"time"
	"github.com/belogik/goes"
	"fmt"
)
func newRemoteWriter() LogInterface {
	nRemoteWriter := &LogRemoteWriter{
		Level: LevelDebug,
	}
	return nRemoteWriter
}

type LogRemoteWriter struct {
	*goes.Connection
	DSN		string 		`json:"Dsn"`
	Level   int			`json:"LogLevel"`
}
// {"dsn":"http://localhost:9200/","levle":1}
func (logRemoteWriter *LogRemoteWriter) Init(jsonConfig string) error {
	err := json.Unmarshal([]byte(jsonConfig), logRemoteWriter)
	if err != nil {
		return nil
	}
	if logRemoteWriter.DSN == "" {
		return errors.New("Empty DSN")
	} else if u, err := url.Parse(logRemoteWriter.DSN); err != nil {
		return err
	} else if u.Path == "" {
		return errors.New("Missing Url Path")
	} else if host, port, err := net.SplitHostPort(u.Host); err != nil {
		return err
	} else {
		conn := goes.NewConnection(host, port)
		logRemoteWriter.Connection = conn
	}
	return nil
}

func (logRemoteWriter *LogRemoteWriter) WriteMessage(when time.Time, message string, level int) error {
	if level >= logRemoteWriter.Level {
		return nil
	}


	r := make(map[string]interface{})
	r["@timestamp"] = when.Format(time.RFC3339)
	r["@msg"] = message
	d := goes.Document{
		Index:	fmt.Sprintf("%04d.%02d.%02d", when.Year(), when.Month(), when.Day()),
		Type: 	"logs",
		Fields: r,
	}
	_, err := logRemoteWriter.Index(d, nil)
	return err
}

func (logRemoteWriter *LogRemoteWriter ) Destroy(){

}

func(logRemoteWriter *LogRemoteWriter) Flush() {

}

func init() {
	Register(AdapterRemote, newRemoteWriter)
}