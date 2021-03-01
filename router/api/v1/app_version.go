package v1

import (
	"net"
	"net/http"
	"time"
	"webapp/appframework"
	"webapp/appframework/app"
	"webapp/appframework/code"
	"webapp/appinterface"
	"webapp/errors"
	"webapp/service"
	"webapp/stat"
	"webapp/toolkit"

	"github.com/gin-gonic/gin"
)

const (
	StatGetAppVersion = "GetAppVersion"
)

func CheckAppVersionApi(c *gin.Context) {
	appframework.BusinessLogger.Infof(c, "content-type:%s", c.Request.Header.Get("Content-Type"))
	t1 := time.Now()
	var form appinterface.ReqBody
	ipSrc := net.ParseIP(c.Request.RemoteAddr)
	payload := int(c.Request.ContentLength)
	_, sessId, err := toolkit.GetUniqId(StatGetAppVersion)
	if err != nil {
		appframework.ErrorLogger.Infof(c, "generate session id failed for:%+v", err)
	}
	err = app.BindAndValid(c, &form)
	appframework.BusinessLogger.Infof(c, "session:%s req body:%+v", sessId, form)
	if err != nil {
		app.JsonResponse(c, http.StatusBadRequest, code.INVALID_PARAMS, err.Error())
		appframework.ErrorLogger.Errorf(c, "session:%s GetAppVersion form: %+v, err: %+v", sessId, form, err)
		go stat.PushStat(StatGetAppVersion, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(code.INVALID_PARAMS))
		return
	}
	result, retCode := service.GetAppVersion(c, sessId, &form.Param.ApiRequest)
	if retCode != errors.RetCode_SUCCESS {
		appframework.ErrorLogger.Errorf(c, "session:%s GetAppVersion form: %+v, err: %+v", sessId, form, err)
		app.JsonResponse(c, http.StatusOK, int(retCode), nil)
		go stat.PushStat(StatGetAppVersion, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(retCode))
		return
	}

	app.JsonResponse(c, http.StatusOK, code.SUCCESS, result)
	go stat.PushStat(StatGetAppVersion, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(errors.RetCode_SUCCESS))
}
