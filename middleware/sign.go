package middleware

import (
	"net/http"
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
			app.JsonResponse(c, http.StatusUnauthorized, code.ERRO_SERVICE_ID_FIELD_NO_EXIST, code.GetMsg(code.ERRO_SERVICE_ID_FIELD_NO_EXIST))
			c.Abort()
			return
		}
		Sign := c.Request.Header.Get(appframework.LocalServiceCfg.CheckSignField)
		if Sign == "" {
			app.JsonResponse(c, http.StatusUnauthorized, code.ERROR_SIGN_FIELD_NO_EXIST, code.GetMsg(code.ERROR_SIGN_FIELD_NO_EXIST))
			c.Abort()
			return
		}
		bIsAllow := toolkit.ArrayCheckIn(callServiceId, appframework.LocalServiceCfg.AllowServiceIdList)
		if !bIsAllow {
			app.JsonResponse(c, http.StatusUnauthorized, code.ERROR_DENY_SERVICE_ID, code.GetMsg(code.ERROR_DENY_SERVICE_ID))
			c.Abort()
			return
		}
		reqData := ""
		for _, v := range appframework.LocalServiceCfg.CheckSignData {
			data := c.Request.Header.Get(v)
			if data == "" {
				app.JsonResponse(c, http.StatusUnauthorized, code.ERROR_LOST_SIGN_DATA, code.GetMsg(code.ERROR_LOST_SIGN_DATA))
				c.Abort()
				return
			}
			reqData = reqData + data
		}
		localSign := toolkit.ApiSign(reqData, appframework.LocalServiceCfg.CheckSignKey)
		if localSign != Sign {
			app.JsonResponse(c, http.StatusUnauthorized, code.ERROR_DENY_SERVICE_ID, code.GetMsg(code.ERROR_DENY_SERVICE_ID))
			c.Abort()
			return
		}

		c.Next()
	}
}
