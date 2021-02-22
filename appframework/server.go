package appframework

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	Version    = "1.0.0"
	AppTypeWeb = 1
)

// Application ...
type Application struct {
	Name           string
	Type           int32
	LoggerRootPath string
	SetupVars      func() error
}

// ListenerApplication ...
type WEBApplication struct {
	*Application
	EndPort        int
	MonitorEndPort int

	// 监控使用的http server
	Mux *http.ServeMux
	// RegisterHttpRoute 定义HTTP router
	RegisterHttpRoute func() *gin.Engine
	// 系统定时任务
	RegisterTasks func() []CronTask
}

type CronTask struct {
	Cron     string
	TaskFunc func()
}