package appinterface

import (
	"webapp/appframework"
	"webapp/toolkit"

	"github.com/zanlichard/beegoe/validation"
)

//基本协议头定义
type ReqHeader struct {
	CallServiceId string `json:"_callerServiceId"`
	GroupNo       string `json:"_groupNo"`
	Interface     string `json:"_interface"`
	InvokeId      string `json:"_invokeId"`
	MsgType       string `json:"_msgType"`
	Remark        string `json:"_remark"`
	Timestamp     string `json:"_timestamps"`
	Version       string `json:"_version"`
}

//检查APP版本请求定义
type AppVersionCheckReq struct {
	ClientType     int8   `valid:"Required" json:"client_type"` //当前版本
	CurrentVersion string `valid:"Required" json:"current_ver"` //客户端端类型(1:ios,2:android,3:web)
}

//检查APP版本响应定义
type AppVersionCheckRsp struct {
	BuildCode   string `json:"build_code"`   // 构建的代码
	DownloadUrl string `json:"download_url"` // 下载的url
	ForceUpdate uint8  `json:"force_update"` // 0否，1是
	VersionName string `json:"version_name"` // 版本名称
	Title       string `json:"title"`        // 标题
	Content     string `json:"content"`      // 内容
	Remark      string `json:"remark"`       // 备注
}

//业务请求嵌套定义
type ParamInfo struct {
	ApiRequest AppVersionCheckReq `valid:"Required" json:"clientinfo"`
}

//基本请求体定义
type ReqBody struct {
	Head  ReqHeader `json:"_head"`
	Param ParamInfo `json:"_param"` //上层应用定义
}

func (t *AppVersionCheckReq) Valid(v *validation.Validation) {
	if t.ClientType != 1 && t.ClientType != 2 {
		v.SetError("ClientType", "ClientType有效期取值只能为1,2")
	}
	if len(t.CurrentVersion) != 6 {
		v.SetError("current_ver", "长度不合法")
	}

}

//基本配置管理接口定义(header)
type BasicCfgGetReq struct {
	CfgType string `valid:"Required" json:"cfg_type"` //rabbitmq,mysql,redis,mongo as the key

}

func (t *BasicCfgGetReq) Valid(v *validation.Validation) {
	supportedTypes := []string{"rabbitmq", "mysql", "redis", "mongo"}
	if !toolkit.ArrayCheckIn(t.CfgType, supportedTypes) {
		v.SetError("cfg_type", "不支持的类型")
	}
}

type BasicCfgGetRsp struct {
	UserName     string   `json:"user_name"`
	Passwd       string   `json:"pass_word"`
	Database     string   `json:"database_name"`
	Hosts        []string `json:"host_names"`
	MaxOpenConns int      `json:"max_open_conns"`
	MaxIdleConns int      `json:"max_idle_conns"`
	IdleTimeout  int      `json:"idle_timeout"`
	MaxActive    int      `json:"max_active"` //redis
	Other        string   `json:"extends"`    //rabbitmq-vhost-queue
}

//依赖配置管理接口定义(header)
type DepCfgGetReq struct {
	IsServicesAll bool   `form:"is_services_all" valid:"Required"` //是否全部读取,false,则需要指定service_name
	ServiceName   string `form:"service_name"`                     //服务名
}

func (t *DepCfgGetReq) Valid(v *validation.Validation) {
	if !t.IsServicesAll {
		if t.ServiceName == "" {
			v.SetError("service_name", "参数不全")
		}

	}
}

type DepCfgGetRsp struct {
	Services []appframework.AclDependentItem `json:"services"`
}

type LocalAclRsp struct {
	LocalCfg appframework.LocalAcl `json:"local_config"`
}
