package v1

import (
	"webapp/pkg/app"
	"webapp/pkg/code"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"webapp/globalconfig"
)

func IndexApi(c *gin.Context) {
	app.JsonResponse(c, http.StatusOK, code.SUCCESS, "Welcome to "+ globalconfig.App.Name)
	return
}

func PingApi(c *gin.Context) {
	app.JsonResponse(c, http.StatusOK, code.SUCCESS, time.Now())
}