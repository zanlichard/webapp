package appframeworkboot

import (
	"webapp/appframework"
	"webapp/router"
	"webapp/globalconfig"
	"context"
	"github.com/gin-gonic/gin"
)

/*
  RegisterHttpRoute 此处注册http接口
  类似nginx的access、error日志
*/
func RegisterHttpRoute() *gin.Engine {
	accessInfoLogger := &AccessInfoLogger{}
	accessErrLogger := &AccessErrLogger{}
	ginRouter := router.InitRouter(accessInfoLogger, accessErrLogger)
	return ginRouter
}

type AccessInfoLogger struct{}

func (a *AccessInfoLogger) Write(p []byte) (n int, err error) {
	globalconfig.AccessLogger.Infof(context.Background(), "[gin-info] %s", p)
	return 0, nil
}

type AccessErrLogger struct{}

func (a *AccessErrLogger) Write(p []byte) (n int, err error) {
	globalconfig.ErrorLogger.Errorf(context.Background(), "[gin-err] %s", p)
	return 0, nil
}

// 注册定时任务
func RegisterTasks() []appframework.CronTask {
	var tasks = make([]appframework.CronTask, 0)
	tasks = append(tasks,
		//TestCronTask(), // 测试定时任务
	)

	return tasks
}
