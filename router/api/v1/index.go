package v1

import (
	"net/http"
	"time"
	"webapp/appframework/app"
	"webapp/appframework/code"

	"webapp/appengine"

	"github.com/gin-gonic/gin"
)

func IndexApi(c *gin.Context) {
	app.JsonResponsev2(c, http.StatusOK, code.SUCCESS, "Welcome to "+appengine.App.Name)
	return
}

func PingApi(c *gin.Context) {
	app.JsonResponsev2(c, http.StatusOK, code.SUCCESS, time.Now())
}
