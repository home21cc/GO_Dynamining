package logs

import (
	"sync"
	"fmt"
	"runtime"
	"path"
	"time"
	"os"
	"strconv"
	"log"
	"strings"
)

const levelLogImpl = -1


// RFC5424 LOG
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

//Legacy Log Level
const (
	LevelInfo 	= LevelInformational
	LevelTrace 	= LevelDebug
	LevelWarn 	= LevelWarning
)

const (
	AdapterConnection	= "AdapterConnection"
	AdapterConsole		= "AdapterConsole"
	AdapterFile			= "AdapterFile"
	AdapterRemote		= "AdapterRemote"
	AdapterSmtp			= "AdapterSmtp"
)

type logInterfaceFunc func() LogInterface

type LogInterface interface {
	Init(config string) error
	WriteMessage(when time.Time, message string, level int) error
	Destroy()
	Flush()
}

var adapters = make(map[string]logInterfaceFunc)

func Register(name string, log logInterfaceFunc) {
	if log == nil {
		panic("Logs: Register provide is nil")
	}

	if _, dup := adapters[name]; dup {
		panic("Logs: Register called twice for provider " + name)
	}
	adapters[name] = log
}

type DynamicLog struct {
	lock 				sync.Mutex
	level				int
	init				bool
	enableFuncCallDepth	bool
	logFuncCallDepth	int
	asynchronous		bool
	msgChanLen			int64
	msgChan				chan *logMessage
	signalChan			chan string
	waitGroup			sync.WaitGroup
	msg					chan *logMessage
	outputs				[]*nameLog
}

const defaultAsyncMessageLen = 1e3

type nameLog struct {
	LogInterface
	name string
}

type logMessage struct {
	level		int
	message		string
	when		time.Time
}

var logMessagePool *sync.Pool

func newLogger(chanLens ...int64) *DynamicLog {
	dynamicLog := new(DynamicLog)
	dynamicLog.level				= LevelDebug
	dynamicLog.logFuncCallDepth 	= 2
	dynamicLog.msgChanLen = append(chanLens, 0)[0]
	if dynamicLog.msgChanLen <=0 {
		dynamicLog.msgChanLen = defaultAsyncMessageLen
	}
	dynamicLog.signalChan = make(chan string, 1)
	dynamicLog.setLogger(AdapterConsole)
	return dynamicLog
}

func (dynamicLog *DynamicLog) Async(messageLen ...int64) *DynamicLog {
	dynamicLog.lock.Lock()
	defer dynamicLog.lock.Unlock()
	if dynamicLog.asynchronous {
		return dynamicLog
	}

	dynamicLog.asynchronous = true
	if len(messageLen) > 0 && messageLen[0] > 0 {
		dynamicLog.msgChanLen = messageLen[0]
	}
	dynamicLog.msgChan = make(chan *logMessage, dynamicLog.msgChanLen)
	logMessagePool = &sync.Pool{
		New: func() interface{} {
			return &logMessage{}
		},
	}
	dynamicLog.waitGroup.Add(1)
	go dynamicLog.startLogger()
	return dynamicLog
}

func (dynamicLog *DynamicLog) setLogger(adapterName string, jsonConfigs ...string) error {
	jsonFile := append(jsonConfigs, "{}")[0]
	for _, l := range dynamicLog.outputs {
		if l.name == adapterName {
			return fmt.Errorf("logs: duplicate adaptername %q (you hanv set this log before", adapterName)
		}
	}

	logAdapter, ok := adapters[adapterName]
	if !ok {
		return fmt.Errorf("logs: unknow adaptername %q (forgetten Register ?", adapterName)
	}

	lg := logAdapter()
	err := lg.Init(jsonFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "logs.DynamicLog.SetDynamicLog: " + err.Error())
		return err
	}

	dynamicLog.outputs = append(dynamicLog.outputs, &nameLog{name: adapterName, LogInterface: lg})
	return nil
}

func (dynamicLog *DynamicLog) SetDynamicLog(adapterName string, jsonFiles ...string) error {
	dynamicLog.lock.Lock()
	defer dynamicLog.lock.Unlock()
	if !dynamicLog.init {
		dynamicLog.outputs = []*nameLog{}
		dynamicLog.init = true
	}
	return dynamicLog.setLogger(adapterName, jsonFiles...)
}

func (dynamicLog *DynamicLog) DeleteDynamicLog(adapterName string) error {
	dynamicLog.lock.Lock()
	defer dynamicLog.lock.Unlock()
	outputs := []*nameLog{}
	for _, lg:= range dynamicLog.outputs {
		if lg.name == adapterName {
			lg.Destroy()
		} else {
			outputs = append(outputs, lg)
		}
	}
	if len(outputs) == len(dynamicLog.outputs) {
		return fmt.Errorf("logs: unknown adapterName %q (forgetten Register ?)", adapterName)
	}
	dynamicLog.outputs = outputs
	return nil
}

func (dynamicLog *DynamicLog) writeToLoggers(when time.Time, message string, level int) {
	for _, l := range dynamicLog.outputs {
		err := l.WriteMessage(when, message, level)
		if err != nil {
			fmt.Fprintf(os.Stderr, "uinable to WriteMessage to adpater: %v, error: %v\n", l.name, err)
		}
	}
}

func (dynamicLog *DynamicLog) Write(p [] byte) (n int, err error) {
	if len(p) == 0 {
		return 0, nil
	}
	if p[len(p)-1] == '\n' {
		p = p[0:len(p)-1]
	}
	err = dynamicLog.writeMessage(levelLogImpl, string(p))
	if err == nil {
		return len(p), err
	}
	return 0, err
}

func (dynamicLog *DynamicLog) writeMessage(logLevel int, message string, v ...interface{}) error {
	if !dynamicLog.init {
		dynamicLog.lock.Lock()
		dynamicLog.setLogger(AdapterFile)
		dynamicLog.lock.Unlock()
	}

	if len(v) > 0 {
		message = fmt.Sprintf(message, v...)
	}
	when := time.Now()

	if dynamicLog.enableFuncCallDepth {
		_, file, line, ok := runtime.Caller(dynamicLog.logFuncCallDepth)
		if !ok {
			file = "???"
			line = 0
		}
		_, filename := path.Split(file)
		message = "[" + filename + ":" + strconv.FormatInt(int64(line), 10) + "] " + message
	}

	if dynamicLog.asynchronous {
		lMessage := logMessagePool.Get().(*logMessage)
		lMessage.level = logLevel
		lMessage.message	= message
		lMessage.when	= when
		dynamicLog.msgChan <- lMessage
	} else {
		dynamicLog.writeToLoggers(when, message, logLevel)
	}
	return nil
}

func (dynamicLog *DynamicLog) SetLevel(level int) {
	dynamicLog.level = level
}

func (dynamicLog *DynamicLog) SetLogFuncCallDepth(depth int) {
	dynamicLog.logFuncCallDepth = depth
}

func (dynamicLog *DynamicLog) GetLogFuncCallDepth() int{
	return dynamicLog.logFuncCallDepth
}

func (dynamicLog *DynamicLog) EnableFuncCallDepth(b bool) {
	dynamicLog.enableFuncCallDepth = b
}

func (dynamicLog *DynamicLog) startLogger() {
	flag := false
	for {
		select {
		case logMessage := <-dynamicLog.msgChan:
			dynamicLog.writeToLoggers(logMessage.when, logMessage.message, logMessage.level)
			logMessagePool.Put(logMessage)
		case signal := <-dynamicLog.signalChan:
			dynamicLog.flush()
			if signal == "close" {
				for _, l := range dynamicLog.outputs {
					l.Destroy()
				}
				dynamicLog.outputs = nil
				flag = true
			}
			dynamicLog.waitGroup.Done()
		}
		if flag {
			break
		}
	}
}

func (dynamicLog *DynamicLog) Emergency(format string, v ...interface{}) {
	if LevelEmergency > dynamicLog.level {
		return
	}
	msg := fmt.Sprintf("[Emergncy] " + format, v...)
	dynamicLog.writeMessage(LevelEmergency, msg, v...)
}

func (dynamicLog *DynamicLog) Alert(format string, v ...interface{}) {
	if LevelAlert > dynamicLog.level {
		return
	}
	msg := fmt.Sprintf("[Alert] " + format, v...)
	dynamicLog.writeMessage(LevelAlert, msg, v...)
}

func (dynamicLog *DynamicLog) Critical(format string, v...interface{}) {
	if LevelCritical > dynamicLog.level {
		return
	}
	msg := fmt.Sprintf("[Critical] " + format, v...)
	dynamicLog.writeMessage(LevelCritical, msg, v...)
}

func (dynamicLog *DynamicLog) Error(format string, v...interface{}) {
	if LevelError > dynamicLog.level {
		return
	}
	msg := fmt.Sprintf("[Error] " + format, v...)
	dynamicLog.writeMessage(LevelError, msg, v...)
}

func (dynamicLog *DynamicLog) Warning(format string, v...interface{}) {
	if LevelWarning > dynamicLog.level {
		return
	}
	msg := fmt.Sprintf("[Warning] " + format, v...)
	dynamicLog.writeMessage(LevelWarning, msg)
}



func (dynamicLog *DynamicLog) Notice(format string, v...interface{}) {
	if LevelNotice > dynamicLog.level {
		return
	}
	msg := fmt.Sprintf("[Notice] " + format, v...)
	dynamicLog.writeMessage(LevelNotice, msg, v...)
}


func (dynamicLog *DynamicLog) Informational(format string, v ...interface{}) {
	if LevelInformational > dynamicLog.level {
		return
	}
	msg := fmt.Sprintf("[Info] " + format, v...)
	dynamicLog.writeMessage(LevelInformational, msg, v...)
}


func (dynamicLog *DynamicLog) Debug(format string, v ...interface{}) {
	if LevelDebug > dynamicLog.level {
		return
	}
	msg := fmt.Sprintf("[Debug] " + format, v...)
	dynamicLog.writeMessage(LevelDebug, msg, v...)
}


func (dynamicLog *DynamicLog) Warn(format string, v ...interface{}) {
	if LevelWarn > dynamicLog.level {
		return
	}
	msg := fmt.Sprintf("[Warning] " + format, v...)
	dynamicLog.writeMessage(LevelWarn, msg, v...)
}


func (dynamicLog *DynamicLog) Info(format string, v ...interface{}) {
	if LevelInfo > dynamicLog.level {
		return
	}
	msg := fmt.Sprintf("[Info] " + format, v...)
	dynamicLog.writeMessage(LevelInfo, msg, v...)
}


func (dynamicLog *DynamicLog) Trace(format string, v ...interface{}) {
	if LevelDebug > dynamicLog.level {
		return
	}
	msg := fmt.Sprintf("[Debug] " + format, v...)
	dynamicLog.writeMessage(LevelDebug, msg, v...)
}

func (dynamicLog *DynamicLog) Flush() {
	if dynamicLog.asynchronous {
		dynamicLog.signalChan <- "flush"
		dynamicLog.waitGroup.Wait()
		dynamicLog.waitGroup.Add(1)
		return
	}
	dynamicLog.flush()
}

func (dynamicLog *DynamicLog) Close() {
	if dynamicLog.asynchronous {
		dynamicLog.signalChan <- "close"
		dynamicLog.waitGroup.Wait()
		close(dynamicLog.msgChan)
	} else {
		dynamicLog.flush()
		for _, l := range dynamicLog.outputs {
			l.Destroy()
		}
		dynamicLog.outputs = nil
	}
	close(dynamicLog.signalChan)
}

func (dynamicLog *DynamicLog) Reset() {
	dynamicLog.Flush()
	for _, l := range dynamicLog.outputs {
		l.Flush()
	}
	dynamicLog.outputs = nil
}

func (dynamicLog *DynamicLog) flush() {
	if dynamicLog.asynchronous {
		for {
			if len(dynamicLog.msgChan) > 0 {
				dMessage := <- dynamicLog.msgChan
				dynamicLog.writeToLoggers(dMessage.when, dMessage.message, dMessage.level)
				logMessagePool.Put(dMessage)
				continue
			}
			break
		}
	}
	for _, l := range dynamicLog.outputs {
		l.Flush()
	}
}



var dynaLog *DynamicLog = newLogger()

func GetDynamicLog() *DynamicLog {
	return dynaLog
}

var dynamicLoggerMap = struct {
	sync.RWMutex
	logs map[string]*log.Logger
}{
	logs: map[string]*log.Logger{},
}

func GetLogger(prefixes ...string) *log.Logger {
	prefix := append(prefixes, "")[0]
	if prefix != "" {
		prefix = fmt.Sprintf(`[%s] `, strings.ToUpper(prefix))
	}
	dynamicLoggerMap.RLock()
	l, ok := dynamicLoggerMap.logs[prefix]
	if ok {
		dynamicLoggerMap.RUnlock()
		return l
	}

	dynamicLoggerMap.RUnlock()
	dynamicLoggerMap.Lock()
	defer dynamicLoggerMap.Unlock()
	l, ok = dynamicLoggerMap.logs[prefix]
	if !ok {
		l = log.New(dynaLog, prefix, 0)
		dynamicLoggerMap.logs[prefix] = l
	}
	return l
}

func Reset() {
	dynaLog.Reset()
}

func Async(msgLen ...int64) *DynamicLog {
	return dynaLog.Async(msgLen...)
}

func SetLevel(l int) {
	dynaLog.SetLevel(l)
}

func EnableFuncCallDepth(b bool) {
	dynaLog.enableFuncCallDepth = b
}

func SetLogFuncCall(b bool) {
	dynaLog.EnableFuncCallDepth(b)
	dynaLog.SetLogFuncCallDepth(4)
}

func SetLogFuncCallDepth(d int) {
	dynaLog.logFuncCallDepth = d
}

func SetLogger(adapter string, configs ...string) error {
	err := dynaLog.SetDynamicLog(adapter, configs...)
	if err != nil {
		return err
	}
	return nil
}

func Emergency(f interface{}, v ...interface{}) {
	dynaLog.Emergency(formatLog(f, v...))
}
func Alert(f interface{}, v ...interface{}) {
	dynaLog.Alert(formatLog(f, v...))
}

func Critical(f interface{}, v ...interface{}) {
	dynaLog.Critical(formatLog(f, v...))
}

func Error(f interface{}, v ...interface{}) {
	dynaLog.Error(formatLog(f, v...))
}

func Notice(f interface{}, v ...interface{}) {
	dynaLog.Notice(formatLog(f, v...))
}

func Warning(f interface{}, v ...interface{}) {
	dynaLog.Warning(formatLog(f, v...))
}

func Warn(f interface{}, v ...interface{}) {
	dynaLog.Warn(formatLog(f, v...))
}

func Informational(f interface{}, v ...interface{}) {
	dynaLog.Info(formatLog(f, v...))
}

func Info(f interface{}, v ...interface{}) {
	dynaLog.Info(formatLog(f, v...))
}

func Debug(f interface{}, v ...interface{}) {
	dynaLog.Debug(formatLog(f, v...))
}

func Trace(f interface{}, v ...interface{}) {
	dynaLog.Trace(formatLog(f, v...))
}

func formatLog(f interface{}, v ...interface{}) string {
	var msg string
	switch f.(type) {
	case string:
		msg = f.(string)
		if len(v) == 0 {
			return msg
		}
		if strings.Contains(msg, "%") && !strings.Contains(msg, "%%") {

		} else {
			msg += strings.Repeat(" %v", len(v))
		}
	default:
		msg = fmt.Sprint(f)
		if len(v) == 0 {
			return msg
		}
		msg += strings.Repeat(" %v", len(v))
	}
	return fmt.Sprintf(msg, v...)
}