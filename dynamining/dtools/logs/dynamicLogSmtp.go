package logs

import (
	"encoding/json"
	"strings"
	"fmt"
	"time"
	"net/smtp"
	"net"
	"crypto/tls"
)

type LogSmtpWriter struct {
	Username			string 		`json:"UserName"`
	Password			string		`json:"Password"`
	Host				string		`json:"Host"`
	Subject				string  	`json:"Subject"`
	FromAddress			string		`json:"FromAddress"`
	RecipientAddress	[]string	`json:"SendTo"`
	Level				int			`json:"LogLevel"`
}

func newSmtpWriter() LogInterface {
	nSmtpWriter := &LogSmtpWriter{
		Level: LevelTrace,
	}
	return nSmtpWriter
}

// JsonConfig file example
//{
//	"UserName": "example@gmail.com",
//	"Password": "password",
//	"Host": "smtp.gmail.com:465",
//	"Subject": "email Title",
//	"FromAddress": "from@gmail.com",
//	"SendTo": ["email1","email2"],
//	"LogLevel": LevelError,
//}
func (logSmtpWriter *LogSmtpWriter) Init(jsonConfig string) error {
	err := json.Unmarshal([]byte(jsonConfig), logSmtpWriter)
	if err != nil {
		return err
	}
	return nil
}

func (logSmtpWriter *LogSmtpWriter) WriteMessage(when time.Time, message string, level int) error {
	if level > logSmtpWriter.Level {
		return nil
	}

	hostPort := strings.Split(logSmtpWriter.Host, ":")
	auth := logSmtpWriter.getSmtpAuth(hostPort[0])
	contentType := "Content-Type: text/plain" + "; charset=UTF-8"
	mailMessage := []byte("To: " + strings.Join(logSmtpWriter.RecipientAddress, ";") + "\r\n" +
		"From: " + logSmtpWriter.FromAddress + "<" + logSmtpWriter.FromAddress + ">" +  "\r\n" +
		"Subject: " + logSmtpWriter.Subject + "\r\n" +	contentType + "\r\n\r\n" +
		fmt.Sprintf(".%s", when.Format("2017-01-01 09:17:55")) + message )
	return logSmtpWriter.sendMail(logSmtpWriter.Host, auth, logSmtpWriter.FromAddress, logSmtpWriter.RecipientAddress, mailMessage)
}

func (logSmtpWriter *LogSmtpWriter) Flush() {
	return
}

func( LogSmtpWriter *LogSmtpWriter) Destroy() {
	return
}

func init() {
	Register(AdapterSmtp, newSmtpWriter)
}

func (logSmtpWriter *LogSmtpWriter) getSmtpAuth(host string) smtp.Auth {
	if len(strings.Trim(logSmtpWriter.Username, " ")) == 0 &&
	len(strings.Trim(logSmtpWriter.Password, " ")) == 0 {
		return nil
	}
	return smtp.PlainAuth(
		"",
		logSmtpWriter.Username,
		logSmtpWriter.Password,
		host,
	)
}

func (logSmtpWriter *LogSmtpWriter) sendMail(hostAddressWithPort string, auth smtp.Auth,
	fromAddress string, recipients []string, messageContent []byte ) error {
	client, err := smtp.Dial(hostAddressWithPort)
	if err != nil {
		return err
	}

	host, _, _ := net.SplitHostPort(hostAddressWithPort)
	tlsConn := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:			host,
	}
	if err = client.StartTLS(tlsConn); err != nil{
		return err
	}
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return err
		}
	}
	if err = client.Mail(fromAddress); err != nil {
		return err
	}
	for _, rec := range recipients {
		if err = client.Rcpt(rec); err != nil {
			return err
		}
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(messageContent))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	err = client.Quit()
	if err != nil {
		return err
	}
	return nil

}