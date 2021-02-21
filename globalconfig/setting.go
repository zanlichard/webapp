package globalconfig

import (
	"webapp/appframework"
	"gitee.com/cristiane/go-common/log"
	"net/http"
	"time"
)


var (
	App                   *appframework.WEBApplication
	AccessLogger          log.LoggerContextIface
	ErrorLogger           log.LoggerContextIface
	BusinessLogger        log.LoggerContextIface
	HttpClient            = &http.Client{Timeout: 30 * time.Second}
)