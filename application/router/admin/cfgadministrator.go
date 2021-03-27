package admin

import (
	"fmt"
	"net"
	"net/http"
	"time"
	"webapp/application/appconfig"
	e "webapp/application/apperrors"
	"webapp/application/appinterface"
	"webapp/frame/appframework"
	"webapp/frame/appframework/app"
	"webapp/frame/appframework/code"
	"webapp/frame/trace"
	"webapp/stat"

	"github.com/gin-gonic/gin"
)

const (
	StatGetBasicCfg     = "GetBasicCfg"
	StatGetDependentCfg = "GetDepCfg"
	StatGetLocalAcl     = "GetLocalAcl"
)

func GetBasicConfig(c *gin.Context) {
	defer trace.Recovery()
	appframework.BusinessLogger.Infof(c, "content-type:%s", c.Request.Header.Get("Content-Type"))
	t1 := time.Now()
	var form appinterface.BasicCfgGetReq
	ipSrc := net.ParseIP(c.Request.RemoteAddr)
	payload := int(c.Request.ContentLength)
	err := app.BindAndValid(c, &form)
	appframework.BusinessLogger.Infof(c, "req body:%+v", form)
	if err != nil {
		app.JsonResponsev2(c, http.StatusBadRequest, code.INVALID_PARAMS, err.Error())
		appframework.ErrorLogger.Errorf(c, "GetBasicCfg form: %+v, err: %+v", form, err)
		go stat.PushStat(StatGetBasicCfg, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(code.INVALID_PARAMS))
		return
	}
	var result appinterface.BasicCfgGetRsp
	switch form.CfgType {
	case "mysql":
		result.UserName = appconfig.Config.Database.Mysql.User
		result.Passwd = appconfig.Config.Database.Mysql.Passwd
		result.Database = appconfig.Config.Database.Mysql.Database
		result.Hosts = []string{appconfig.Config.Database.Mysql.ServerAddr}
		result.MaxOpenConns = appconfig.Config.Database.Mysql.MaxOpenConns
		result.MaxIdleConns = appconfig.Config.Database.Mysql.MaxIdleConns
		result.IdleTimeout = appconfig.Config.Database.Mysql.IdleTimeout

	case "redis":
		result.Passwd = appconfig.Config.RedisInfo.Passwd
		result.Hosts = appconfig.Config.RedisInfo.ServerList
		result.MaxOpenConns = appconfig.Config.RedisInfo.MaxIdle
		result.IdleTimeout = appconfig.Config.RedisInfo.IdleTimeout
		result.MaxActive = appconfig.Config.RedisInfo.MaxActive

	case "rabbitmq":
		result.UserName = appconfig.Config.RabbitMq.Username
		result.Passwd = appconfig.Config.RabbitMq.Password
		result.Hosts = []string{appconfig.Config.RabbitMq.ServerAddr}
		result.Other = "{\"queuename\":apptoml.Config.RabbitMq.Queuename,\"vhost\":apptoml.Config.RabbitMq.Vhost}"

	case "mongo":
		result.UserName = appconfig.Config.Mongodb.Username
		result.Passwd = appconfig.Config.Mongodb.Password
		result.Database = appconfig.Config.Mongodb.DB
		host := fmt.Sprintf("%s:%d", appconfig.Config.Mongodb.Server, appconfig.Config.Mongodb.Port)
		result.Hosts = []string{host}

	default:
		result.UserName = ""
		result.Passwd = ""
		result.Database = ""
		result.MaxOpenConns = 0
		result.MaxIdleConns = 0
		result.IdleTimeout = 0
	}
	app.JsonResponsev2(c, http.StatusOK, code.SUCCESS, result)
	go stat.PushStat(StatGetBasicCfg, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(e.RetCode_SUCCESS))
}

func GetLocalAclConfig(c *gin.Context) {
	defer trace.Recovery()
	appframework.BusinessLogger.Infof(c, "content-type:%s", c.Request.Header.Get("Content-Type"))
	t1 := time.Now()
	ipSrc := net.ParseIP(c.Request.RemoteAddr)
	payload := int(c.Request.ContentLength)
	var result appinterface.LocalAclRsp
	result.LocalCfg = appframework.LocalServiceCfg
	app.JsonResponsev2(c, http.StatusOK, code.SUCCESS, result)
	go stat.PushStat(StatGetLocalAcl, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(e.RetCode_SUCCESS))

}

func GetDependentConfig(c *gin.Context) {
	defer trace.Recovery()
	appframework.BusinessLogger.Infof(c, "content-type:%s", c.Request.Header.Get("Content-Type"))
	t1 := time.Now()
	var form appinterface.DepCfgGetReq
	ipSrc := net.ParseIP(c.Request.RemoteAddr)
	payload := int(c.Request.ContentLength)
	err := app.BindAndValid(c, &form)
	appframework.BusinessLogger.Infof(c, "req body:%+v", form)
	if err != nil {
		app.JsonResponsev2(c, http.StatusBadRequest, code.INVALID_PARAMS, err.Error())
		appframework.ErrorLogger.Errorf(c, "GetDepCfg form: %+v, err: %+v", form, err)
		go stat.PushStat(StatGetDependentCfg, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(code.INVALID_PARAMS))
		return
	}
	var result appinterface.DepCfgGetRsp
	for _, item := range appframework.DependenceServiceMap {
		result.Services = append(result.Services, item)
	}
	app.JsonResponsev2(c, http.StatusOK, code.SUCCESS, result)
	go stat.PushStat(StatGetDependentCfg, int(time.Now().Sub(t1).Seconds()*1000), ipSrc, payload, int(e.RetCode_SUCCESS))

}

func SetBasicConfig(c *gin.Context) {
	defer trace.Recovery()
	return
}
