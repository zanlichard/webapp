package router

import (
	"io"
	"os"
	"webapp/apptoml"
	"webapp/middleware"
	v1 "webapp/router/api/v1"
	"webapp/stat"

	"github.com/gin-gonic/gin"
)

func initStat() {
	stat.GStat.AddReportBodyRowItem(v1.StatGetAppVersion)
	stat.GStat.AddReportErrorItem(v1.StatGetAppVersion)
}

func InitRouter(accessInfoLogger, accessErrLogger io.Writer) *gin.Engine {
	gin.DefaultWriter = io.MultiWriter(os.Stdout, accessInfoLogger)
	gin.DefaultErrorWriter = io.MultiWriter(os.Stderr, accessErrLogger)

	if apptoml.Config.Server.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.Cors())
	r.Use(middleware.CheckCallSign())

	r.GET("/", v1.IndexApi)
	r.GET("/ping", v1.PingApi)
	apiG := r.Group("/api")
	apiV1 := apiG.Group("/v1")
	{
		apiResources := apiV1.Group("/app")
		{
			apiResources.POST("/check_version", v1.CheckAppVersionApi)
		}
	}
	//初始化监控
	initStat()
	return r
}
