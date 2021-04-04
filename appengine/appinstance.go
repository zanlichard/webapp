package appengine

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"webapp/application/appconfig"
	"webapp/application/apperrors"
	"webapp/application/router"
	"webapp/application/storage"
	"webapp/frame/appframework"
	"webapp/logger"
	"webapp/stat"
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
	logger.InfoFormat("level:%s path:%s NamePrefix:%s filename:%s interval:%d",
		appconfig.Config.Server.Stat.LogLevel,
		appconfig.Config.Server.Stat.LogPath,
		appconfig.Config.Server.Stat.NamePrefix,
		appconfig.Config.Server.Stat.Filename,
		appconfig.Config.Server.Stat.Interval)
	stat.Init(*logConfig, time.Duration(appconfig.Config.Server.Stat.Interval))
	stat.SetDelayUp(50, 100, 200)
	stat.Proc()
}

func initMongo() bool {
	mongoHost := fmt.Sprintf("%s:%d", appconfig.Config.Mongodb.Server, appconfig.Config.Mongodb.Port)
	err := storage.InitMgo(mongoHost, appconfig.Config.Mongodb.DB, appconfig.Config.Mongodb.Username, appconfig.Config.Mongodb.Password, 20, false)
	if err != nil {
		logger.ErrorFormat("init mongo failed for:%+v ", err)
		return false
	}
	return true
}

func initDB() bool {
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
		return false
	}
	return true
}

func RegisterHttpRoute(isDebug bool) *gin.Engine {
	accessInfoLogger := &appframework.AccessInfoLogger{}
	accessErrLogger := &appframework.AccessErrLogger{}
	ginRouter := router.InitRouter(accessInfoLogger, accessErrLogger, isDebug)
	logger.InfoFormat("init router begin")
	return ginRouter
}

func ReLoad(serviceName string) {

}

func InitAppInstance(serviceName string) bool {
	//初始化配置系统
	appconfig.Init("")
	//初始化本地调用栈跟踪
	apperrors.Init(serviceName)
	//初始化日志系统
	initLogger(serviceName)
	//初始化本地监控
	initStat()
	logger.InfoFormat("init stat finish")
	//初始化数据库
	if !initDB() {
		logger.ErrorFormat("init database failed")
		return false
	}
	logger.InfoFormat("init database finish")
	//初始化mongodb
	if !initMongo() {
		logger.ErrorFormat("init mongodb failed")
		return false
	}
	logger.InfoFormat("init mongodb finish")
	//加载服务依赖
	if err := appframework.InitServiceDependence(serviceName, appconfig.Config.ConfigMng.DepServiceList); err != nil {
		logger.ErrorFormat("init service dependence cfg err:%+v", err.Error())
		return false
	}
	logger.InfoFormat("init service dependence finish")
	//加载本地访问访问规则
	if err := appframework.InitServiceLocalCfg(appconfig.Config.ConfigMng.AclServiceList); err != nil {
		logger.ErrorFormat("init service local acl cfg err:%+v", err.Error())
		return false
	}
	logger.InfoFormat("init local acl finish")
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
