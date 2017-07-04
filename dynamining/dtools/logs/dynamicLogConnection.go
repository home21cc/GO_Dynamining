package logs

import (
	"io"
	"net"
	"encoding/json"
	"time"
)

type ConnWriter struct {
	logConnAdapter		*logWriter
	innerWriter			io.WriteCloser
	ReconnectOnMessage	bool	`json:"ReconnectOnMessage"`
	Reconnect			bool	`json:"Reconnect"`
	Net					string	`json:"Net"`
	Address				string 	`json:"Addr"`
	Level				int		`json:"LogLevel"`
}

func newConn() LogInterface {
	conn := new(ConnWriter)
	conn.Level = LevelTrace
	return conn
}

// jsonConfig only need key "Level"
func (connWriter *ConnWriter) Init(jsonConfig string) error {
	return json.Unmarshal([]byte(jsonConfig), connWriter)
}

func (connWriter *ConnWriter) WriteMessage(when time.Time, message string, level int) error {
	if level > connWriter.Level {
		return nil
	}

	if connWriter.neededConnectOnMsg() {
		err := connWriter.connect()
		if err != nil {
			return err
		}
	}

	if connWriter.ReconnectOnMessage {
		defer connWriter.innerWriter.Close()
	}
	connWriter.logConnAdapter.println(when, message)
	return nil
}

func (connWriter *ConnWriter) Flush() {

}

func (connWriter *ConnWriter) Destroy() {
	if connWriter.innerWriter != nil {
		connWriter.innerWriter.Close()
	}
}

func (connWriter *ConnWriter) connect() error {
	if connWriter.innerWriter != nil {
		connWriter.innerWriter.Close()
		connWriter.innerWriter = nil
	}

	conn, err := net.Dial(connWriter.Net, connWriter.Address)
	if err != nil {
		return err
	}

	if tcpConn, ok := conn.(*net.TCPConn); ok {
		tcpConn.SetKeepAlive(true)
	}

	connWriter.innerWriter = conn
	connWriter.logConnAdapter = newLogWriter(conn)
	return nil
}

func (connWriter *ConnWriter) neededConnectOnMsg() bool {
	if connWriter.Reconnect {
		connWriter.Reconnect = false
		return true
	}

	if connWriter.innerWriter == nil {
		return true
	}
	return connWriter.ReconnectOnMessage
}

func init() {
	Register(AdapterConnection, newConn)
}