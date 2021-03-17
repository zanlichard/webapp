package router

import (
	"io"
	"os"
	"webapp/apptoml"
	"webapp/middleware"
	"webapp/router/admin"
	v1 "webapp/router/api/v1"
	"webapp/router/resmgr"
	"webapp/stat"

	"github.com/gin-contrib/pprof"

	"github.com/gin-gonic/gin"
)

func initStat() {
	stat.GStat.AddReportBodyRowItem(v1.StatGetAppVersion)
	stat.GStat.AddReportBodyRowItem(admin.StatGetBasicCfg)
	stat.GStat.AddReportBodyRowItem(admin.StatGetDependentCfg)
	stat.GStat.AddReportBodyRowItem(admin.StatGetLocalAcl)

	stat.GStat.AddReportErrorItem(v1.StatGetAppVersion)
	stat.GStat.AddReportErrorItem(admin.StatGetBasicCfg)
	stat.GStat.AddReportErrorItem(admin.StatGetDependentCfg)
	stat.GStat.AddReportErrorItem(admin.StatGetLocalAcl)
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

	if apptoml.Config.Server.Debug {
		pprof.Register(r)
	}

	r.Use(middleware.Cors())
	r.GET("/", v1.IndexApi)
	r.GET("/ping", v1.PingApi)
	apiAdmin := r.Group("/admin")
	{
		apiAdmin.POST("/get-basic-cfg", admin.GetBasicConfig)
		apiAdmin.POST("/get-dep-cfg", admin.GetDependentConfig)
		apiAdmin.POST("/set-basic-cfg", admin.SetBasicConfig)
		apiAdmin.POST("/get-local-acl", admin.GetLocalAclConfig)
	}

	apiResMng := r.Group("/resourcemgr")
	{
		apiResMng.POST("/heartbeat", resmgr.Heartbeat)
	}
	apiG := r.Group("/api")
	r.Use(middleware.CheckCallSign())
	//r.Use(middleware.CheckUserToken())
	apiV1 := apiG.Group("/v1")
	{
		apiResources := apiV1.Group("/app")
		{
			apiResources.POST("/check-version", v1.CheckAppVersionApi)
		}
	}
	//初始化监控
	initStat()
	return r
}
