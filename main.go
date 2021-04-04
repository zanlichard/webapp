package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"webapp/appengine"
	"webapp/logger"
)

const (
	//Version 版本
	Version = "010000"
	//VersionEx 版本
	VersionEx = "1.0.0"
	//Update 版本
	Update = "2021-2-19 17:46:00"
	//服务名
	AppName = "webapp"
)

func signHandler() {
	c := make(chan os.Signal)
	signal.Notify(c,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGUSR1,

		syscall.SIGUSR2)
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Printf("recv signal:%+v", s)
				appengine.ExitAppInstance()
				os.Exit(-1)
			case syscall.SIGUSR1:
			case syscall.SIGUSR2:
				fmt.Println("reload config")
				appengine.ReLoad(AppName)
			default:
				fmt.Printf("other signal:%+v", s)
			}
		}
	}()

}

func initEnv() {
	//设置允许调度运行的CPU个数
	runtime.GOMAXPROCS(runtime.NumCPU())
	//补充信号处理
	signHandler()
}

func main() {
	//运行环境初始化
	initEnv()

	defer appengine.ExitAppInstance()
	//应用实例初始化
	if !appengine.InitAppInstance(AppName) {
		fmt.Printf("%s init application instance failed", logger.GetLogDatePrefix())
		return
	}

	//应用实例启动
	fmt.Printf("%s begin to start appInstance", logger.GetLogDatePrefix())
	appengine.StartAppInstance()

}
