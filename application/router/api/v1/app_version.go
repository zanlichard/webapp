package v1

import (
	"net"
	"net/http"
	"time"
	e "webapp/application/apperrors"
	"webapp/application/appinterface"
	"webapp/application/service"
	"webapp/frame/appframework"
	"webapp/frame/appframework/app"
	"webapp/frame/appframework/code"
	"webapp/frame/trace"
	"webapp/stat"
	"webapp/toolkit"

	"github.com/gin-gonic/gin"
)

const (
	StatGetAppVersion = "GetAppVersion"
)

func CheckAppVersionApi(c *gin.Context) {
	defer trace.Recovery()
	appframework.BusinessLogger.Infof(c, "content-type:%s", c.Request.Header.Get("Content-Type"))
	t1 := time.Now()
	var reqMsg appinterface.AppVerCheckMsg
	ipSrc := net.ParseIP(c.Request.RemoteAddr)
	payload := int(c.Request.ContentLength)
	_, sessId, err := toolkit.GetUniqId(StatGetAppVersion)
	if err != nil {
		appframework.BusinessLogger.Infof(c, "generate session id failed for:%+v", err)
	}
	rspHeader := reqMsg.Head
	rspHeader.MsgType = "response"
	rspHeader.Timestamp = toolkit.ConvertToString(toolkit.GetTimeStamp())

	err = app.BindOnly(c, &reqMsg)
	if err != nil {
		app.JsonResponse(c, http.StatusBadRequest, code.ERROR_BAD_REQUEST, rspHeader, nil)
		appframework.BusinessLogger.Errorf(c, "session:%s GetAppVersion err: %+v", sessId, err)
		go stat.PushStat(StatGetAppVersion, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(code.INVALID_PARAMS))
		return
	}

	err = app.ValidOnly(c, &reqMsg.Param)
	appframework.BusinessLogger.Infof(c, "session:%s req body:%+v", sessId, reqMsg.Param)
	if err != nil {
		app.JsonResponse(c, http.StatusBadRequest, code.INVALID_PARAMS, rspHeader, nil)
		appframework.BusinessLogger.Errorf(c, "session:%s GetAppVersion form: %+v, err: %+v", sessId, reqMsg.Param, err)
		go stat.PushStat(StatGetAppVersion, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(code.INVALID_PARAMS))
		return
	}
	result, retCode := service.GetAppVersion(c, sessId, &reqMsg.Param)
	if retCode != e.RetCode_SUCCESS {
		appframework.BusinessLogger.Errorf(c, "session:%s GetAppVersion form: %+v, err: %+v", sessId, reqMsg.Param, err)
		app.JsonResponse(c, http.StatusOK, int(retCode), rspHeader, nil)
		go stat.PushStat(StatGetAppVersion, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(retCode))
		return
	}

	app.JsonResponse(c, http.StatusOK, code.SUCCESS, rspHeader, result)
	go stat.PushStat(StatGetAppVersion, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(e.RetCode_SUCCESS))
}
