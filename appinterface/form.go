package appinterface

import (
	"github.com/zanlichard/beegoe/validation"
)

//检查APP版本
type AppVersionCheckReq struct {
	ClientType     int8   `form:"client_type" valid:"Required" json:"client_type"` //当前版本
	CurrentVersion string `form:"client_type" valid:"Required" json:"current_ver"` //客户端端类型(1:ios,2:android,3:web)
}
type AppVersionCheckRsp struct {
	BuildCode   string `json:"build_code"`    // 构建的代码
	DownloadUrl string `json:"download_url"`  // 下载的url
	ForceUpdate uint8  `json:"force_update"`  // 0否，1是
	VersionName string `json:"version_name"`  // 版本名称
	Title       string `json:"title"`         // 标题
	Content     string `json:"content"`       // 内容
	Remark      string `json:"remark"`        // 备注
}


func (t *AppVersionCheckReq) Valid(v *validation.Validation) {
	if t.ClientType != 1 && t.ClientType != 2 {
		 v.SetError("ClientType", "ClientType有效期取值只能为1,2")
	}
	if len(t.CurrentVersion) != 6 {
		 v.SetError("current_ver","长度不合法")
	}

}

type Page struct {
	PageSize int `form:"page_size" json:"page_size"`
	PageNum  int `form:"page_num"  json:"page_num"`
}









