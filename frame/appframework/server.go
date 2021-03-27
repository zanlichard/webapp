package appframework

import (
	"net/http"

	"gitee.com/cristiane/go-common/log"

	"github.com/gin-gonic/gin"
)

const (
	Version        = "1.0.0"
	AppTypeWeb     = 1
	LocalServiceId = "2160037"
)

// Application ...
type Application struct {
	Name           string
	Type           int32
	LoggerRootPath string
	LoggerLevel    string
	SetupVars      func() error
}

// ListenerApplication ...
type WEBApplication struct {
	*Application
	EndPort        int
	MonitorEndPort int
	IsDebug        bool
	// 监控使用的http server
	Mux *http.ServeMux
	// RegisterHttpRoute 定义HTTP router
	RegisterHttpRoute func(isDebug bool) *gin.Engine
	// 系统定时任务
	RegisterTasks func() []CronTask
}

type CronTask struct {
	Cron     string
	TaskFunc func()
}

var (
	AccessLogger   log.LoggerContextIface
	ErrorLogger    log.LoggerContextIface
	BusinessLogger log.LoggerContextIface
)

// SetupVars 加载变量
func SetupVars() error {
	var err error
	ErrorLogger, err = log.GetErrLogger("err")
	if err != nil {
		return err
	}

	AccessLogger, err = log.GetAccessLogger("access")
	if err != nil {
		return err
	}

	BusinessLogger, err = log.GetBusinessLogger("business")
	if err != nil {
		return err
	}
	return nil
}
