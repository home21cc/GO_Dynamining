package logs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"strconv"
)

type LogFileWriter struct {
	sync.RWMutex

	Filename 			string 	`json:"Filename"`
	LogFileAdapter 		*os.File

	MaxLines			int		`json:"MaxLines"`
	MaxLinesCurLines	int

	MaxSize				int		`json:"MaxSize"`
	MaxSizeCurSize		int

	Daily				bool 	`json:"Daily"`
	MaxDays				int64	`json:"MaxDays"`
	dailyOpenDate		int
	dailyOpenTime		time.Time

	Rotate				bool	`json:"Rotate"`

	Level				int 	`json:"LogLevel"`

	Perm 				string 	`json:"Perm"`

	fileNameOnly, suffix 	string
}

func newFileWriter() LogInterface {
	nfWriter := &LogFileWriter{
		Filename:	"",
		MaxLines:	1000000,
		MaxSize:	1 << 28,		// 256mb
		Daily:		true,
		MaxDays:	7,
		Rotate:		true,
		Level:		LevelTrace,
		Perm:		"0660",
	}
	return nfWriter
}

func (logWriter *LogFileWriter) Init(jsonConfig string) error {
	err := json.Unmarshal([]byte(jsonConfig), logWriter)
	if err != nil {
		return err
	}

	if len(logWriter.Filename) == 0 {
		return errors.New("jsonConfig must have filename")
	}

	logWriter.suffix = filepath.Ext(logWriter.Filename)
	logWriter.fileNameOnly = strings.TrimSuffix(logWriter.Filename, logWriter.suffix)
	if logWriter.suffix == "" {
		logWriter.suffix = ".log"
	}

	err = logWriter.startLog()
	return err
}

func (logWriter *LogFileWriter) startLog() error {
	file, err := logWriter.createLogFile()
	if err != nil {
		return err
	}
	if logWriter.LogFileAdapter != nil {
		logWriter.LogFileAdapter.Close()
	}
	logWriter.LogFileAdapter = file

	return logWriter.initFile()
}

func (logWriter *LogFileWriter) needRotate(size, day int) bool {
	return  (logWriter.MaxLines > 0  && logWriter.MaxLinesCurLines >= logWriter.MaxLines) ||
	(logWriter.MaxSize > 0 && logWriter.MaxSizeCurSize >= logWriter.MaxSize) ||
	(logWriter.Daily && day != logWriter.dailyOpenDate)

}

func (logWriter *LogFileWriter) WriteMessage(when time.Time, message string, level int) error {
	if level > logWriter.Level {
		return nil
	}
	head, detail := formatTimeHeader(when)
	message = string(head) + message + "\n"
	if logWriter.Rotate {
		logWriter.RLock()
		if logWriter.needRotate(len(message), detail) {
			logWriter.RUnlock()
			logWriter.Lock()
			if logWriter.needRotate(len(message), detail) {
				if err := logWriter.doRotate(when); err != nil {
					fmt.Fprintf(os.Stderr, "LogFileWriter(%q) : %s\n", logWriter.Filename, err )
				}
			}
			logWriter.Unlock()
		} else {
			logWriter.RUnlock()
		}
	}
	logWriter.Lock()
	_, err := logWriter.LogFileAdapter.Write([]byte(message))
	if err == nil {
		logWriter.MaxLinesCurLines++
		logWriter.MaxSizeCurSize += len(message)
	}
	logWriter.Unlock()
	return err
}

func (logWriter *LogFileWriter) createLogFile() (*os.File, error) {
	perm, err := strconv.ParseInt(logWriter.Perm, 8, 64)
	if err != nil {
		return nil, err
	}
	file, err := os.OpenFile(logWriter.Filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.FileMode(perm))
	if err == nil {
		os.Chmod(logWriter.Filename, os.FileMode(perm))
	}
	return file, err
}

func (logWriter *LogFileWriter) initFile() error {
	file := logWriter.LogFileAdapter
	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("get stat err: %s\n", err)
	}
	logWriter.MaxSizeCurSize = int(fileInfo.Size())
	logWriter.dailyOpenTime = time.Now()
	logWriter.dailyOpenDate = logWriter.dailyOpenTime.Day()
	logWriter.MaxLinesCurLines = 0

	if logWriter.Daily {
		go logWriter.dailyRotate(logWriter.dailyOpenTime)
	}
	if fileInfo.Size() > 0  {
		count, err := logWriter.lines()
		if err != nil {
			return err
		}
		logWriter.MaxLinesCurLines = count
	}
	return nil
}

func(logWriter *LogFileWriter) dailyRotate(openTime time.Time) {
	year, month, day := openTime.Add(24 * time.Hour).Date()
	nextDay := time.Date(year, month, day, 0, 0, 0, 0, openTime.Location())
	tm := time.NewTimer(time.Duration(nextDay.UnixNano()-openTime.UnixNano() + 100))
	select {
	case <- tm.C:
		logWriter.Lock()
		if logWriter.needRotate(0, time.Now().Day()) {
			if err := logWriter.doRotate(time.Now()); err != nil {
				fmt.Fprintf(os.Stderr, "LogFileWriter(%q): %s\n", logWriter.Filename, err)
			}
		}
		logWriter.Unlock()
	}
}

func (logWriter *LogFileWriter) lines() (int, error) {
	file, err := os.Open(logWriter.Filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	buf := make([]byte, 32768)	// 32k
	count := 0
	lineSeprate := []byte{'\n'}

	for {
		c, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}

		count += bytes.Count(buf[:c], lineSeprate)
		if err == io.EOF {
			break
		}
	}
	return count, nil
}

func (logWriter *LogFileWriter) doRotate(logTime time.Time) error {
	num := 1
	fileName := ""
	_, err := os.Lstat(logWriter.Filename)
	if err != nil {
		goto RESTART_LOGGER
	}
	if logWriter.MaxLines > 0 || logWriter.MaxSize > 0 {
		for ; err == nil && num <= 999; num++ {
			fileName = logWriter.fileNameOnly + fmt.Sprintf(".%s.%03d%s", logTime.Format("2017-03-22"), num, logWriter.suffix)
			_, err = os.Lstat(fileName)
		}
	} else {
		fileName = fmt.Sprintf("%s.%s%s", logWriter.fileNameOnly, logWriter.dailyOpenTime.Format("2017-03-30"), logWriter.suffix)
		_, err = os.Lstat(fileName)
		for ; err == nil && num <= 999; num++ {
			fileName = logWriter.fileNameOnly + fmt.Sprintf(".%s.%03d%s", logWriter.dailyOpenTime.Format("2017-03-30"), num, logWriter.suffix)
			_, err = os.Lstat(fileName)
		}
	}

	if err == nil {
		return fmt.Errorf("Rotate: Cannot find free log number to rename %s\n", logWriter.Filename)
	}
	logWriter.LogFileAdapter.Close()
	err = os.Rename(logWriter.Filename, fileName)
	err = os.Chmod(fileName, os.FileMode(440))

	RESTART_LOGGER:
	startLoggerErr := logWriter.startLog()
	go logWriter.deleteOLDLog()

	if startLoggerErr != nil {
		return fmt.Errorf("Rotate startLog: %s\n", startLoggerErr)
	}

	if err != nil {
		return fmt.Errorf("Rotate: %s\n", err)
	}
	return nil
}

func (logWriter *LogFileWriter) deleteOLDLog() {
	dir := filepath.Dir(logWriter.Filename)
	filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) (rError error) {
			defer func() {
				if r := recover(); r != nil {
					rError = fmt.Errorf("Unable to delete old log '%s', error: %v\n", path, r)
					fmt.Println(rError)
				}
			}()
			if !info.IsDir() && info.ModTime().Add(24*time.Hour*time.Duration(logWriter.MaxDays)).Before(time.Now()) {
				if strings.HasPrefix(filepath.Base(path), filepath.Base(logWriter.fileNameOnly)) &&
					strings.HasPrefix(filepath.Base(path), logWriter.suffix) {
					os.Remove(path)
				}
			}
		return
	})
}

func (logWriter *LogFileWriter) Destroy() {
	logWriter.LogFileAdapter.Close()
}

func (LogWriter *LogFileWriter) Flush() {
	LogWriter.LogFileAdapter.Sync()
}

func init() {
	Register(AdapterFile, newFileWriter)
}