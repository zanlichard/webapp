package main

import (
	"webapp/storage"
	"webapp/apptoml"
	"webapp/logger"
	"webapp/stat"
	"webapp/appengine"
	"webapp/appframeworkboot"
	"webapp/appframework"
	"fmt"
	"runtime"
	"time"
)
const (
	//Version 版本
	Version = "010000"
	//VersionEx 版本
	VersionEx = "1.0.0"
	//Update 版本
	Update = "2021-2-19 17:46:00"
    //服务名
	AppName = "webapp"
)

func initLogger() {
	if len(apptoml.Config.Server.Log.LogFile) > 0 {
		filename := fmt.Sprintf("%s/%s",
			apptoml.Config.Server.Log.LogDir, apptoml.Config.Server.Log.LogFile)
		config := logger.NewConfig(filename,
			apptoml.Config.Server.Log.MaxLines, apptoml.Config.Server.Log.MaxSize, apptoml.Config.Server.Log.MaxDays)
		logger.Logger = logger.NewLogger("file",
			apptoml.Config.Server.Log.LogLevel, config, apptoml.Config.Server.Log.ChanLen)
	}
	if len(apptoml.Config.Server.Log.AnalysisFile) > 0 {
		filename := fmt.Sprintf("%s/%s",
			apptoml.Config.Server.Log.LogDir, apptoml.Config.Server.Log.AnalysisFile)
		config := logger.NewConfig(filename,
			apptoml.Config.Server.Log.MaxLines, apptoml.Config.Server.Log.MaxSize, apptoml.Config.Server.Log.MaxDays)
		logger.AnalysisLogger = logger.NewLogger("file",
			apptoml.Config.Server.Log.LogLevel, config, apptoml.Config.Server.Log.ChanLen)
	}
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	initLogger()
}

func exitEnv() {
	logger.CloseLogger()
}

func initStat() {
	logconfig := new(stat.LoggerParam)
	logconfig.Level = apptoml.Config.Server.Stat.LogLevel
	logconfig.Path = apptoml.Config.Server.Stat.LogPath
	logconfig.NamePrefix = apptoml.Config.Server.Stat.NamePrefix
	logconfig.Filename = apptoml.Config.Server.Stat.Filename
	logconfig.Maxfilesize = apptoml.Config.Server.Stat.MaxFileSize
	logconfig.Maxdays = apptoml.Config.Server.Stat.MaxDays
	logconfig.Maxlines = apptoml.Config.Server.Stat.MaxLines
	fmt.Printf("level:%s path:%s NamePrefix:%s filename:%s interval:%d\n",
		apptoml.Config.Server.Stat.LogLevel,
		apptoml.Config.Server.Stat.LogPath,
		apptoml.Config.Server.Stat.NamePrefix,
		apptoml.Config.Server.Stat.Filename,
		apptoml.Config.Server.Stat.Interval)
	stat.Init(*logconfig, time.Duration(apptoml.Config.Server.Stat.Interval))
	stat.SetDelayUp(50, 100, 200)
	stat.Proc()

}

func exitStat() {
	stat.Exit()

}

func main() {
	//应用层初始化
	initEnv()
	defer exitEnv()

	//本地化监控初始化
	initStat()
	defer exitStat()

	//存储初始化
	serverAddr := apptoml.Config.Database.Mysql.ServerAddr
	user       := apptoml.Config.Database.Mysql.User
	pwd        := apptoml.Config.Database.Mysql.Passwd
	dbase      := apptoml.Config.Database.Mysql.Database
	maxOpen    := apptoml.Config.Database.Mysql.MaxOpenConns
	maxIdle    := apptoml.Config.Database.Mysql.MaxIdleConns
	idleTime   := apptoml.Config.Database.Mysql.IdleTimeout
	storage.InitDB(serverAddr,user,pwd,dbase,maxOpen,maxIdle,idleTime)

	//框架初始化
	application := &appframework.WEBApplication{
		Application: &appframework.Application{
			Name:       AppName,
			SetupVars:  appframeworkboot.SetupVars,
		},
		RegisterHttpRoute:appframeworkboot.RegisterHttpRoute,
		RegisterTasks:appframeworkboot.RegisterTasks,
	}
	appengine.RunApplication(application)

}
