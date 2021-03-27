package v1

import (
	"net/http"
	"time"
	"webapp/frame/appframework/app"
	"webapp/frame/appframework/code"

	"github.com/gin-gonic/gin"
)

func IndexApi(c *gin.Context) {
	app.JsonResponsev2(c, http.StatusOK, code.SUCCESS, "Welcome to webapp")
	return
}

func PingApi(c *gin.Context) {
	app.JsonResponsev2(c, http.StatusOK, code.SUCCESS, time.Now())
}
