package appengine

import (
	"webapp/appframework"
	"webapp/apptoml"
	"webapp/internal/setup"
	. "webapp/logger"

	"gitee.com/cristiane/go-common/log"
)

// 初始化application--日志部分
func initApplication(application *appframework.Application) error {
	err := log.InitGlobalConfig(apptoml.Config.Server.Log.LogDir, apptoml.Config.Server.Log.LogLevel, application.Name)
	if err != nil {
		Logger.Error("InitGlobalConfig:%+v", err)
		return err
	}
	return nil
}

// 初始化监控相关的http接口
func setupCommonVars(application *appframework.WEBApplication) error {
	if apptoml.Config.Server.MonitorEndPort != 0 {
		application.Mux = setup.NewServerMux()
	}
	return nil
}
