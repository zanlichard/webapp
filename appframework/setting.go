package appframework

type CronTask struct {
	Cron string
	TaskFunc func()
}