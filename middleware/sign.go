package middleware

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"webapp/appframework"
	"webapp/appframework/app"
	"webapp/appframework/code"
	"webapp/toolkit"

	"github.com/gin-gonic/gin"
)

func CheckCallSign() gin.HandlerFunc {
	return func(c *gin.Context) {
		callServiceId := c.Request.Header.Get(appframework.LocalServiceCfg.CheckIdField)
		if callServiceId == "" {
			app.JsonResponse(c, http.StatusUnauthorized, code.ERRO_SERVICE_ID_FIELD_NO_EXIST, "")
			appframework.ErrorLogger.Errorf(c, "get headerfield:%s failed", appframework.LocalServiceCfg.CheckIdField)
			c.Abort()
			return
		}
		Sign := c.Request.Header.Get(appframework.LocalServiceCfg.CheckSignField)
		if Sign == "" {
			app.JsonResponse(c, http.StatusUnauthorized, code.ERROR_SIGN_FIELD_NO_EXIST, "")
			appframework.ErrorLogger.Errorf(c, "get header field:%s failed", appframework.LocalServiceCfg.CheckSignField)
			c.Abort()
			return
		}
		bIsAllow := toolkit.ArrayCheckIn(callServiceId, appframework.LocalServiceCfg.AllowServiceIdList)
		if !bIsAllow {
			app.JsonResponse(c, http.StatusUnauthorized, code.ERROR_DENY_SERVICE_ID, "")
			appframework.ErrorLogger.Errorf(c, "get header field:%s failed", strings.Join(appframework.LocalServiceCfg.AllowServiceIdList, ", "))
			c.Abort()
			return
		}

		body, _ := ioutil.ReadAll(c.Request.Body)
		reqData := string(body)
		localSign := toolkit.ApiSign(reqData, appframework.LocalServiceCfg.CheckSignKey)
		if localSign != Sign {
			app.JsonResponse(c, http.StatusUnauthorized, code.ERROR_DENY_SERVICE_ID, "")
			appframework.ErrorLogger.Errorf(c, "request sign:%s local sign:%s", Sign, localSign)
			c.Abort()
			return
		}
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		c.Next()
	}
}
