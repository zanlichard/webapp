package appengine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"webapp/application/appconfig"
	"webapp/application/apperrors"
	"webapp/application/router"
	"webapp/frame/appframework"
	"webapp/logger"
	"webapp/stat"
	"webapp/storage"
)

var (
	App *appframework.WEBApplication
)

func initLogger(serviceName string) {
	if len(appconfig.Config.Server.Log.LogFile) > 0 {
		filename := fmt.Sprintf("%s/%s",
			appconfig.Config.Server.Log.LogDir, appconfig.Config.Server.Log.LogFile)
		config := logger.NewConfig(filename,
			appconfig.Config.Server.Log.MaxLines, appconfig.Config.Server.Log.MaxSize, appconfig.Config.Server.Log.MaxDays)
		logger.Logger = logger.NewLogger(serviceName, "file",
			appconfig.Config.Server.Log.LogLevel, config, appconfig.Config.Server.Log.ChanLen)
	}
}

func initStat() {
	logConfig := new(stat.LoggerParam)
	logConfig.Level = appconfig.Config.Server.Stat.LogLevel
	logConfig.Path = appconfig.Config.Server.Stat.LogPath
	logConfig.NamePrefix = appconfig.Config.Server.Stat.NamePrefix
	logConfig.Filename = appconfig.Config.Server.Stat.Filename
	logConfig.Maxfilesize = appconfig.Config.Server.Stat.MaxFileSize
	logConfig.Maxdays = appconfig.Config.Server.Stat.MaxDays
	logConfig.Maxlines = appconfig.Config.Server.Stat.MaxLines
	fmt.Printf("level:%s path:%s NamePrefix:%s filename:%s interval:%d\n",
		appconfig.Config.Server.Stat.LogLevel,
		appconfig.Config.Server.Stat.LogPath,
		appconfig.Config.Server.Stat.NamePrefix,
		appconfig.Config.Server.Stat.Filename,
		appconfig.Config.Server.Stat.Interval)
	stat.Init(*logConfig, time.Duration(appconfig.Config.Server.Stat.Interval))
	stat.SetDelayUp(50, 100, 200)
	stat.Proc()
}

func initDB() {
	//存储初始化
	serverAddr := appconfig.Config.Database.Mysql.ServerAddr
	user := appconfig.Config.Database.Mysql.User
	pwd := appconfig.Config.Database.Mysql.Passwd
	dbase := appconfig.Config.Database.Mysql.Database
	maxOpen := appconfig.Config.Database.Mysql.MaxOpenConns
	maxIdle := appconfig.Config.Database.Mysql.MaxIdleConns
	idleTime := appconfig.Config.Database.Mysql.IdleTimeout
	debug := appconfig.Config.Server.Debug
	err := storage.InitDB(serverAddr, user, pwd, dbase, maxOpen, maxIdle, idleTime, debug)
	if err != nil {
		logger.ErrorFormat("init database err:%+v", err.Error())
		return
	}
}

func RegisterHttpRoute(isDebug bool) *gin.Engine {
	accessInfoLogger := &appframework.AccessInfoLogger{}
	accessErrLogger := &appframework.AccessErrLogger{}
	ginRouter := router.InitRouter(accessInfoLogger, accessErrLogger, isDebug)
	return ginRouter
}

func InitAppInstance(serviceName string) bool {
	//初始化配置系统
	appconfig.Init("")
	//初始化本地调用栈跟踪
	apperrors.Init(serviceName)
	//初始化日志系统
	initLogger(serviceName)
	//加载服务依赖
	if err1 := appframework.InitServiceDependence(serviceName, appconfig.Config.ConfigMng.DepServiceList); err1 != nil {
		logger.ErrorFormat("init service cfg err:%+v", err1.Error())
		return false
	}
	//框架初始化
	App = &appframework.WEBApplication{
		Application: &appframework.Application{
			Name:      serviceName,
			SetupVars: appframework.SetupVars,
		},
		RegisterHttpRoute: RegisterHttpRoute,
		RegisterTasks:     appframework.RegisterTasks,
	}
	appframework.InitApplication(App, serviceName,
		appconfig.Config.Server.Debug,
		appconfig.Config.Server.EndPort,
		appconfig.Config.Server.MonitorEndPort,
		appconfig.Config.Server.Log.FrameworkLog,
		appconfig.Config.Server.Log.LogLevel)
	return true
}

func StartAppInstance() {
	appframework.RunApplication(App)
}

func ExitAppInstance() {
	stat.Exit()
	storage.ExitDB()
	logger.CloseLogger()
}
