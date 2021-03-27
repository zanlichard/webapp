package logger

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"runtime"
	"strconv"
	"time"
	"webapp/toolkit"

	"github.com/zanlichard/beegoe/logs"
)

const (
	ctLayout = "20060102-150405"
)

var (
	Logger     *logs.BeeLogger
	webappName string
)

func logLevel(lv string) int {
	level := logs.LevelDebug
	switch lv {
	case "info":
		level = logs.LevelInformational
	case "notice":
		level = logs.LevelNotice
	case "warn":
		level = logs.LevelWarning
	case "error":
		level = logs.LevelError
	case "critical":
		level = logs.LevelCritical
	case "alert":
		level = logs.LevelAlert
	case "emergency":
		level = logs.LevelEmergency
	}
	return level
}

func NewConfig(filename string, maxlines, maxsize int64, maxdays int) string {
	return fmt.Sprintf(`{
                  "filename":"%s",
                  "maxlines":%d,
                  "maxsize":%d,
                  "maxdays":%d,
                  "perm": "664",
				  "rotateperm":"444",
                  "blankprefix":true
				}`,
		filename, maxlines, maxsize, maxdays)
}

func NewLogger(appName, adapterName, level, config string, channelLens int64) *logs.BeeLogger {
	webappName = appName
	l := logs.NewLogger(channelLens)
	lv := logLevel(level)
	l.SetLogger(adapterName, config)
	l.SetLevel(lv)
	l.BlankPrefix()
	l.Async()
	return l
}

func CloseLogger() {
	Logger.Flush()
	Logger.Close()
}

/*
  提前格式化,保留少部分的数据交给日志库去控制
  将日志的级别、进程ID、协程ID也放置到本级打印

*/
func getLogHost() string {
	ipAddr := ""
	ip, err := toolkit.GetLocalIp()
	if err != nil {
		ipAddr = "127.0.0.1"
	} else {
		ipAddr = ip.String()
	}
	return ipAddr
}

func GetGoroutineID() uint64 {
	b := make([]byte, 64)
	runtime.Stack(b, false)
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

func ErrorFormat(format string, v ...interface{}) {
	var appLog string
	if len(v) > 0 {
		appLog = fmt.Sprintf(format, v...)
	} else {
		appLog = format
	}
	funcName, fileName, lineNo := GetContextInfo()
	dateStr := GetLogDatePrefix()
	ProcId := os.Getpid()
	ThreadId := toolkit.ConvertToString(GetGoroutineID())
	LogLevel := "ERROR"
	IpAddr := getLogHost()
	logData := fmt.Sprintf("%s|%s|%s|%s|%d|%s|%s|%s|%s|%d|%s", dateStr, LogLevel, ThreadId, ThreadId, ProcId, IpAddr, webappName, fileName, funcName, lineNo, appLog)
	Logger.Error("%s", logData)
}

func InfoFormat(format string, v ...interface{}) {
	var appLog string
	if len(v) > 0 {
		appLog = fmt.Sprintf(format, v...)
	} else {
		appLog = format
	}
	funcName, fileName, lineNo := GetContextInfo()
	dateStr := GetLogDatePrefix()
	ProcId := os.Getpid()
	ThreadId := toolkit.ConvertToString(GetGoroutineID())
	LogLevel := "INFO"
	IpAddr := getLogHost()
	logData := fmt.Sprintf("%s|%s|%s|%s|%d|%s|%s|%s|%s|%d|%s", dateStr, LogLevel, ThreadId, ThreadId, ProcId, IpAddr, webappName, fileName, funcName, lineNo, appLog)
	Logger.Info("%s", logData)
}

func DebugFormat(format string, v ...interface{}) {
	var appLog string
	if len(v) > 0 {
		appLog = fmt.Sprintf(format, v...)
	} else {
		appLog = format
	}
	funcName, fileName, lineNo := GetContextInfo()
	dateStr := GetLogDatePrefix()
	ProcId := os.Getpid()
	ThreadId := toolkit.ConvertToString(GetGoroutineID())
	LogLevel := "DEBUG"
	IpAddr := getLogHost()
	logData := fmt.Sprintf("%s|%s|%s|%s|%d|%s|%s|%s|%s|%d|%s", dateStr, LogLevel, ThreadId, ThreadId, ProcId, IpAddr, webappName, fileName, funcName, lineNo, appLog)
	Logger.Debug("%s", logData)
}

func GetLogDatePrefix() string {
	currentTime := time.Now()
	milliSecond := currentTime.UTC().UnixNano() / int64(time.Millisecond)
	leftMs := milliSecond % 1000
	dateStr := currentTime.Format(ctLayout)
	return fmt.Sprintf("%s-%d", dateStr, leftMs)
}

func GetContextInfo() (string, string, int) {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	_, filename := path.Split(file)
	funcName := runtime.FuncForPC(pc).Name()
	return funcName, filename, line
}
