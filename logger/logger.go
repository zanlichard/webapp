package logger

import (
	"fmt"

	"github.com/zanlichard/beegoe/logs"
)

var (
	Logger         *logs.BeeLogger
	AnalysisLogger *logs.BeeLogger
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
				  "rotateperm":"444"
				}`,
		filename, maxlines, maxsize, maxdays)
}

func NewLogger(adapterName, level, config string, channelLens int64) *logs.BeeLogger {
	l := logs.NewLogger(channelLens)
	lv := logLevel(level)
	l.SetLogger(adapterName, config)
	l.SetLevel(lv)
	l.Async()
	return l
}

func CloseLogger() {
	Logger.Flush()
	Logger.Close()

	AnalysisLogger.Flush()
	AnalysisLogger.Close()
}
