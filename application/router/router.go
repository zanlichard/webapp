package router

import (
	"github.com/gin-contrib/pprof"
	"io"
	"os"
	admin2 "webapp/application/router/admin"
	v12 "webapp/application/router/api/v1"
	resmgr2 "webapp/application/router/resmgr"
	"webapp/frame/middleware"
	"webapp/stat"

	"github.com/gin-gonic/gin"
)

func initStat() {
	stat.GStat.AddReportBodyRowItem(v12.StatGetAppVersion)
	stat.GStat.AddReportBodyRowItem(admin2.StatGetBasicCfg)
	stat.GStat.AddReportBodyRowItem(admin2.StatGetDependentCfg)
	stat.GStat.AddReportBodyRowItem(admin2.StatGetLocalAcl)

	stat.GStat.AddReportErrorItem(v12.StatGetAppVersion)
	stat.GStat.AddReportErrorItem(admin2.StatGetBasicCfg)
	stat.GStat.AddReportErrorItem(admin2.StatGetDependentCfg)
	stat.GStat.AddReportErrorItem(admin2.StatGetLocalAcl)
}

func InitRouter(accessInfoLogger, accessErrLogger io.Writer, isDebug bool) *gin.Engine {
	gin.DefaultWriter = io.MultiWriter(os.Stdout, accessInfoLogger)
	gin.DefaultErrorWriter = io.MultiWriter(os.Stderr, accessErrLogger)

	if isDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	if isDebug {
		pprof.Register(r)
	}

	r.Use(middleware.Cors())
	r.GET("/", v12.IndexApi)
	r.GET("/ping", v12.PingApi)
	apiAdmin := r.Group("/admin")
	{
		apiAdmin.POST("/get-basic-cfg", admin2.GetBasicConfig)
		apiAdmin.POST("/get-dep-cfg", admin2.GetDependentConfig)
		apiAdmin.POST("/set-basic-cfg", admin2.SetBasicConfig)
		apiAdmin.POST("/get-local-acl", admin2.GetLocalAclConfig)
	}

	apiResMng := r.Group("/resourcemgr")
	{
		apiResMng.POST("/heartbeat", resmgr2.Heartbeat)
	}
	apiG := r.Group("/api")
	r.Use(middleware.CheckCallSign())
	//r.Use(middleware.CheckUserToken())
	apiV1 := apiG.Group("/v1")
	{
		apiResources := apiV1.Group("/app")
		{
			apiResources.POST("/check-version", v12.CheckAppVersionApi)
		}
	}
	//初始化监控
	initStat()
	return r
}
