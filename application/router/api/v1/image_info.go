package v1

import (
	"github.com/gin-gonic/gin"
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
)

const (
	StatGetImage = "GetImageInfo"
)

func GetImageInfo(c *gin.Context) {
	defer trace.Recovery()
	appframework.BusinessLogger.Infof(c, "content-type:%s", c.Request.Header.Get("Content-Type"))
	t1 := time.Now()
	var form appinterface.GetImageMsg
	ipSrc := net.ParseIP(c.Request.RemoteAddr)
	payload := int(c.Request.ContentLength)
	_, sessId, err := toolkit.GetUniqId(StatGetAppVersion)
	if err != nil {
		appframework.ErrorLogger.Infof(c, "generate session id failed for:%+v", err)
	}
	err = app.BindOnly(c, &form)
	rspHeader := form.Head
	rspHeader.MsgType = "response"
	rspHeader.Timestamp = toolkit.ConvertToString(toolkit.GetTimeStamp())
	if err != nil {
		app.JsonResponse(c, http.StatusBadRequest, code.INVALID_PARAMS, rspHeader, nil)
		appframework.ErrorLogger.Errorf(c, "session:%s GetImageInfo form: %+v, err: %+v", sessId, form, err)
		go stat.PushStat(StatGetImage, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(code.INVALID_PARAMS))
		return
	}
	err = app.ValidOnly(c, form.Param)
	appframework.BusinessLogger.Infof(c, "session:%s req body:%+v", sessId, form)
	if err != nil {
		app.JsonResponse(c, http.StatusBadRequest, code.INVALID_PARAMS, rspHeader, nil)
		appframework.ErrorLogger.Errorf(c, "session:%s GetImageInfo form: %+v, err: %+v", sessId, form, err)
		go stat.PushStat(StatGetImage, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(code.INVALID_PARAMS))
		return
	}
	req, ok := (form.Param).(appinterface.GetImageReq)
	if !ok {
		app.JsonResponse(c, http.StatusBadRequest, code.INVALID_PARAMS, rspHeader, nil)
		appframework.ErrorLogger.Errorf(c, "session:%s GetImageInfo form: %+v param invalid", sessId, form)
		go stat.PushStat(StatGetImage, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(code.INVALID_PARAMS))
		return
	}

	result, retCode := service.GetImageInfo(c, sessId, &req)
	if retCode != e.RetCode_SUCCESS {
		appframework.ErrorLogger.Errorf(c, "session:%s GetImageInfo form: %+v, err: %+v", sessId, form, err)
		app.JsonResponse(c, http.StatusOK, int(retCode), rspHeader, nil)
		go stat.PushStat(StatGetImage, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(retCode))
		return
	}

	app.JsonResponse(c, http.StatusOK, code.SUCCESS, rspHeader, result)
	go stat.PushStat(StatGetImage, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(e.RetCode_SUCCESS))
}
