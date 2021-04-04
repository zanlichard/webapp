package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"webapp/frame/appframework"
	"webapp/frame/appframework/app"
	"webapp/frame/appframework/code"
	"webapp/frame/subsys"
	"webapp/toolkit"

	"github.com/gin-gonic/gin"
)

func CheckCallSign() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqMsg := subsys.SubsysReqBody{}
		body, _ := ioutil.ReadAll(c.Request.Body)
		err := json.Unmarshal(body, &reqMsg)
		if err != nil {
			app.JsonResponse(c, http.StatusBadRequest, code.ERRO_SERVICE_ID_FIELD_NO_EXIST, subsys.SubsysGetBadHeader(), nil)
			appframework.BusinessLogger.Errorf(c, "get headerfield:%s failed", appframework.LocalServiceCfg.CheckIdField)
			c.Abort()
			return
		}

		rspHead := reqMsg.Head
		rspHead.MsgType = "response"
		rspHead.Timestamp = toolkit.ConvertToString(toolkit.GetTimeStamp())

		callServiceId := c.Request.Header.Get(appframework.LocalServiceCfg.CheckIdField)
		if callServiceId == "" {
			app.JsonResponse(c, http.StatusUnauthorized, code.ERRO_SERVICE_ID_FIELD_NO_EXIST, rspHead, nil)
			appframework.BusinessLogger.Errorf(c, "get headerfield:%s failed", appframework.LocalServiceCfg.CheckIdField)
			c.Abort()
			return
		}
		Sign := c.Request.Header.Get(appframework.LocalServiceCfg.CheckSignField)
		if Sign == "" {
			app.JsonResponse(c, http.StatusUnauthorized, code.ERROR_SIGN_FIELD_NO_EXIST, rspHead, nil)
			appframework.BusinessLogger.Errorf(c, "get header field:%s failed", appframework.LocalServiceCfg.CheckSignField)
			c.Abort()
			return
		}
		bIsAllow := toolkit.ArrayCheckIn(callServiceId, appframework.LocalServiceCfg.AllowServiceIdList)
		if !bIsAllow {
			app.JsonResponse(c, http.StatusUnauthorized, code.ERROR_DENY_SERVICE_ID, rspHead, nil)
			appframework.BusinessLogger.Errorf(c, "get header field:%s failed", strings.Join(appframework.LocalServiceCfg.AllowServiceIdList, ", "))
			c.Abort()
			return
		}

		reqData := string(body)
		serviceName := appframework.AclServiceId2Name[callServiceId]
		aclItem := appframework.AclServiceMap[serviceName]
		localSign := toolkit.ApiSign(reqData, aclItem.ServiceKey)
		if localSign != Sign {
			app.JsonResponse(c, http.StatusUnauthorized, code.ERROR_DENY_SERVICE_ID, rspHead, nil)
			appframework.BusinessLogger.Errorf(c, "request sign:%s local sign:%s", Sign, localSign)
			c.Abort()
			return
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		c.Next()
	}
}
