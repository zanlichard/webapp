package appframeworkboot

import (
	"webapp/appframework"
	"log"
)

const (
	Cron = "0 */10 * * * *"
)

func TestCronTask() appframework.CronTask {
	var task = appframework.CronTask{
		Cron:Cron,
		TaskFunc: func() {
			log.Println("需要一些定时任务支持抽奖环节")
		},
	}
	return task
}
