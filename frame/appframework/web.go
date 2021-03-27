package appframework

import (
	"fmt"
	"gitee.com/cristiane/go-common/log"
	"net/http"
	"strconv"
	"webapp/frame/internal/setup"
	"webapp/logger"
	"webapp/xml"
)

var (
	LocalServiceCfg          LocalAcl
	ServicenameDependenceMap map[string]AclDependentItem //servicename as the key
	ServiceIdDependenceMap   map[string]string           //serviceId map to serviceName
)

//访问控制配置选项
type AclDependentItem struct {
	ServiceName string `json:"service_name"`
	ServiceId   string `json:"service_id"`
	ServiceKey  string `json:"service_key"`
	ServiceAlg  string `json:"service_sign_algorithm"`
	ServiceUrl  string `json:"service_url"`
}

//依赖的外部服务的配置记录集
type DownstreamDependence struct {
	ServiceDependence []AclDependentItem `json:"dependent_declare"`
}

//服务本地配置
type LocalAcl struct {
	AllowServiceIdList []string `json:"allow_ids"`
	CheckIdField       string   `json:"check_id_field"`
	CheckSignField     string   `json:"check_sign_field"`
	CheckSignData      []string `json:"check_sign_data"` //_head,_param
	CheckSignKey       string   `json:"check_key"`
	CheckAlgorithm     string   `json:"check_algorithm"`
	LocalServiceId     string   `json:"local_id"`
}

//从配置管理中心获取服务的自身配置和依赖的外部服务的配置
func InitServiceDependence(serviceName string, dependenceServices []string) error {
	configFile := "./etc/dependence.xml"
	err := xml.XmlInit(configFile)
	if err != nil {
		return err
	}
	defer xml.XmlFree()
	m := make(map[string]AclDependentItem)

	for _, service := range dependenceServices {
		logger.InfoFormat("service:%s", service)
		url := fmt.Sprintf("/DEPENDENTSERVERINFO/%s/Url", service)
		id := fmt.Sprintf("/DEPENDENTSERVERINFO/%s/Id", service)
		name := fmt.Sprintf("/DEPENDENTSERVERINFO/%s/Name", service)
		key := fmt.Sprintf("/DEPENDENTSERVERINFO/%s/Key", service)

		serviceUrl := xml.XmlGetField(url)
		serviceId := xml.XmlGetField(id)
		serviceName := xml.XmlGetField(name)
		serviceKey := xml.XmlGetField(key)
		acl := AclDependentItem{
			ServiceName: serviceName,
			ServiceId:   serviceId,
			ServiceKey:  serviceKey,
			ServiceAlg:  "md5",
			ServiceUrl:  serviceUrl,
		}
		m[serviceName] = acl
	}
	ServicenameDependenceMap = m
	//id2name := make(map[string]string)
	//ServiceIdDependenceMap = id2name
	return nil

}

//初始化本地访问控制
func initServiceLocalCfg() {
	acl := LocalAcl{
		LocalServiceId:     "2160037",
		CheckAlgorithm:     "md5",
		CheckSignKey:       "h6F2GvOm1Q1pR5ATYbMjUIUyscLiBs3E",
		AllowServiceIdList: []string{"2120013", "2160034"},
		CheckIdField:       "HSB-OPENAPI-CALLERSERVICEID",
		CheckSignField:     "HSB-OPENAPI-SIGNATURE",
		CheckSignData:      []string{"_head", "_param"},
	}
	LocalServiceCfg = acl
}
func InitApplication(app *WEBApplication, appName string, isDebug bool, endPoint int, monitorEndPoint int, logPath string, logLevel string) {
	app.EndPort = endPoint
	app.LoggerRootPath = logPath
	app.Type = AppTypeWeb
	app.MonitorEndPort = monitorEndPoint
	app.LoggerLevel = logLevel
	app.Name = appName
	app.IsDebug = isDebug
	initServiceLocalCfg()
	return
}

func RunApplication(app *WEBApplication) {
	if app.Name == "" {
		logger.ErrorFormat("Application name can't not be empty")
		return
	}
	if app.LoggerLevel == "" {
		logger.ErrorFormat("Application loglevel can't not be empty")
		return
	}
	if app.LoggerRootPath == "" {
		logger.ErrorFormat("Application log path can't not be empty")
		return
	}

	err := runApp(app)
	if err != nil {
		logger.ErrorFormat("App.RunListenerApplication err: %v", err)
	}
}

func runApp(webApp *WEBApplication) error {
	err := log.InitGlobalConfig(webApp.LoggerRootPath, webApp.LoggerLevel, webApp.Name)
	if err != nil {
		logger.ErrorFormat("InitGlobalConfig:%+v", err)
		return err
	}
	return nil

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
		logger.InfoFormat("App run monitor server addr: %v", addr)
		err := http.ListenAndServe(addr, webApp.Mux)
		if err != nil {
			logger.ErrorFormat("App run monitor server err: %v", err)
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
	}

	// 7. run http server
	if webApp.RegisterHttpRoute == nil {
		logger.ErrorFormat("App RegisterHttpRoute nil ??")
	}
	// 8. register and gin framework startup
	err = webApp.RegisterHttpRoute(webApp.IsDebug).Run(addr)

	return err
}

// 初始化监控相关的http接口
func setupCommonVars(application *WEBApplication) error {
	application.Mux = setup.NewServerMux()
	return nil
}

// setupGRPCVars ...
func setupWEBVars(webApp *WEBApplication) error {
	err := setupCommonVars(webApp)
	if err != nil {
		return err
	}
	return nil
}
