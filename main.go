package main

import (
	"fmt"
	"runtime"
	"time"
	"webapp/appengine"
	"webapp/apperrors"
	"webapp/appframework"
	"webapp/appframeworkboot"
	"webapp/apptoml"
	"webapp/logger"
	"webapp/stat"
	"webapp/storage"
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
		logger.Logger = logger.NewLogger(AppName, "file",
			apptoml.Config.Server.Log.LogLevel, config, apptoml.Config.Server.Log.ChanLen)
	}
}

func initEnv() {
	//初始化配置系统
	apptoml.Init("")

	apperrors.Init(AppName)

	runtime.GOMAXPROCS(runtime.NumCPU())

	//初始化日志系统
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

func releaseDB() {
	storage.ExitDB()
}

//从配置管理中心获取服务的自身配置和依赖的外部服务的配置
func initServiceDependence() error {
	configFile := "./etc/dependence.xml"
	err := apptoml.ServiceDepInit(configFile)
	if err != nil {
		return err
	}
	defer apptoml.ServiceDepFree()

	m := make(map[string]appframework.AclDependentItem)

	for _, service := range apptoml.Config.ConfigMng.DepServiceList {
		logger.InfoFormat("service:%s", service)
		url := fmt.Sprintf("/DEPENDENTSERVERINFO/%s/Url", service)
		id := fmt.Sprintf("/DEPENDENTSERVERINFO/%s/Id", service)
		name := fmt.Sprintf("/DEPENDENTSERVERINFO/%s/Name", service)
		key := fmt.Sprintf("/DEPENDENTSERVERINFO/%s/Key", service)

		serviceUrl := apptoml.ServiceDepGetField(url)
		serviceId := apptoml.ServiceDepGetField(id)
		serviceName := apptoml.ServiceDepGetField(name)
		serviceKey := apptoml.ServiceDepGetField(key)
		acl := appframework.AclDependentItem{
			ServiceName: serviceName,
			ServiceId:   serviceId,
			ServiceKey:  serviceKey,
			ServiceAlg:  "md5",
			ServiceUrl:  serviceUrl,
		}
		m[serviceName] = acl
	}
	appframework.ServicenameDependenceMap = m
	acl := appframework.LocalAcl{
		LocalServiceId:     "2160037",
		CheckAlgorithm:     "md5",
		CheckSignKey:       "h6F2GvOm1Q1pR5ATYbMjUIUyscLiBs3E",
		AllowServiceIdList: []string{"2120013", "2160034"},
		CheckIdField:       "HSB-OPENAPI-CALLERSERVICEID",
		CheckSignField:     "HSB-OPENAPI-SIGNATURE",
		CheckSignData:      []string{"_head", "_param"},
	}

	appframework.LocalServiceCfg = acl
	localServiceId := "2160037"
	id2name := make(map[string]string)
	id2name[localServiceId] = AppName
	appframework.ServiceIdDependenceMap = id2name
	return nil

}

func main() {
	//运行环境初始化
	initEnv()
	defer exitEnv()

	//本地化监控初始化
	initStat()
	defer exitStat()

	//加载服务依赖
	if err1 := initServiceDependence(); err1 != nil {
		logger.ErrorFormat("init service cfg err:%+v", err1.Error())
		return
	}

	//存储初始化
	serverAddr := apptoml.Config.Database.Mysql.ServerAddr
	user := apptoml.Config.Database.Mysql.User
	pwd := apptoml.Config.Database.Mysql.Passwd
	dbase := apptoml.Config.Database.Mysql.Database
	maxOpen := apptoml.Config.Database.Mysql.MaxOpenConns
	maxIdle := apptoml.Config.Database.Mysql.MaxIdleConns
	idleTime := apptoml.Config.Database.Mysql.IdleTimeout
	debug := apptoml.Config.Server.Debug
	err := storage.InitDB(serverAddr, user, pwd, dbase, maxOpen, maxIdle, idleTime, debug)
	if err != nil {
		logger.ErrorFormat("init database err:%+v", err.Error())
		return
	}
	defer releaseDB()

	//框架初始化
	application := &appframework.WEBApplication{
		Application: &appframework.Application{
			Name:      AppName,
			SetupVars: appframeworkboot.SetupVars,
		},
		RegisterHttpRoute: appframeworkboot.RegisterHttpRoute,
		RegisterTasks:     appframeworkboot.RegisterTasks,
	}
	appengine.RunApplication(application)

}
