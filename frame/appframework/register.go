package appframework

import (
	"context"
)

type AccessInfoLogger struct{}

func (a *AccessInfoLogger) Write(p []byte) (n int, err error) {
	AccessLogger.Infof(context.Background(), "[gin-info] %s", p)
	return 0, nil
}

type AccessErrLogger struct{}

func (a *AccessErrLogger) Write(p []byte) (n int, err error) {
	ErrorLogger.Errorf(context.Background(), "[gin-err] %s", p)
	return 0, nil
}

// 注册定时任务
func RegisterTasks() []CronTask {
	var tasks = make([]CronTask, 0)
	tasks = append(tasks) //TestCronTask(), // 测试定时任务
	return tasks
}
