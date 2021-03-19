package appengine

import (
	"fmt"
	"net/http"
	"strconv"
	"webapp/appframework"
	"webapp/apptoml"
	. "webapp/logger"
)

var (
	App *appframework.WEBApplication
)

func RunApplication(application *appframework.WEBApplication) {
	if application.Name == "" {
		ErrorFormat("Application name can't not be empty")
	}

	application.EndPort = apptoml.Config.Server.EndPort
	application.LoggerRootPath = apptoml.Config.Server.Log.LogDir
	application.Type = appframework.AppTypeWeb
	application.MonitorEndPort = apptoml.Config.Server.MonitorEndPort

	err := runApp(application)
	if err != nil {
		ErrorFormat("App.RunListenerApplication err: %v", err)
	}
	App = application
}

func runApp(webApp *appframework.WEBApplication) error {
	// 1. init application
	err := initApplication(webApp.Application)
	if err != nil {
		return err
	}

	// 2. setup vars
	err = setupWEBVars(webApp)
	if err != nil {
		return err
	}
	if webApp.SetupVars != nil {
		err = webApp.SetupVars()
		if err != nil {
			return fmt.Errorf("App.SetupVars err: %v", err)
		}
	}

	//3.  setup server monitor in single goroutine
	go func() {
		addr := "0.0.0.0:" + strconv.Itoa(webApp.MonitorEndPort)
		InfoFormat("App run monitor server addr: %v", addr)
		err := http.ListenAndServe(addr, webApp.Mux)
		if err != nil {
			ErrorFormat("App run monitor server err: %v", err)
		}
	}()

	// 5 run task
	//cn := cron.New(cron.WithSeconds())
	//cronTasks := webApp.RegisterTasks()
	//for i := 0;i<len(cronTasks);i++{
	//	if cronTasks[i].TaskFunc != nil {
	//		_,err = cn.AddFunc(cronTasks[i].Cron,cronTasks[i].TaskFunc)
	//		if err != nil {
	//			logging.Fatalf("App run cron task err: %v",err)
	//		}
	//	}
	//}
	//cn.Start()

	// 6. set init service port
	var addr string
	if webApp.EndPort != 0 {
		addr = "0.0.0.0:" + strconv.Itoa(webApp.EndPort)
	} else if apptoml.Config.Server.EndPort != 0 {
		addr = "0.0.0.0:" + strconv.Itoa(apptoml.Config.Server.EndPort)
	}

	// 7. run http server
	if webApp.RegisterHttpRoute == nil {
		ErrorFormat("App RegisterHttpRoute nil ??")
	}
	// 8. register and gin framework startup
	err = webApp.RegisterHttpRoute().Run(addr)

	return err
}

// setupGRPCVars ...
func setupWEBVars(webApp *appframework.WEBApplication) error {
	err := setupCommonVars(webApp)
	if err != nil {
		return err
	}

	return nil
}
