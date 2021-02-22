package appframework

import (
	"net/http"

	"gitee.com/cristiane/go-common/log"

	"github.com/gin-gonic/gin"
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

//访问控制配置选项
type AclDependentItem struct {
	ServiceName string `json:"service_name"`
	ServiceId   string `json:"service_id"`
	ServiceKey  string `json:"service_key"`
	ServiceAlg  string `json:"service_sign_algorithm"`
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

var (
	AccessLogger             log.LoggerContextIface
	ErrorLogger              log.LoggerContextIface
	BusinessLogger           log.LoggerContextIface
	ServicenameDependenceMap map[string]AclDependentItem //servicename as the key
	LocalServiceCfg          LocalAcl
	ServiceIdDependenceMap   map[string]string //serviceId map to serviceName
)
