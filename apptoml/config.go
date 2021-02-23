package apptoml

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

var (
	configFile = os.Getenv("CONFIG") //"./config.toml"
)

var (
	Config = globalConfig{}
)

type (
	globalConfig struct {
		Title     string    `toml:"title"`
		Server    server    `toml:"server"`
		Database  database  `toml:"database"`
		Redisinfo redisinfo `toml:"redisinfo"`
		RabbitMq  rabbitmq  `toml:"rabbitmq"`
		ConfigMng cfgcenter `toml:"cfgcenter"`
	}
	rabbitmq struct {
		Username   string `toml:"username"`
		Password   string `toml:"password"`
		ServerAddr string `toml:"server"`
		ServerPort int    `toml:"port"`
		Vhost      string `toml:"vhost"`
		Queuename  string `toml:"queue"`
	}

	server struct {
		Debug              bool   `toml:"debug"`
		ServiceName        string `toml:"serviceName"`
		Log                log    `toml:"log"`
		Stat               stat   `toml:"stat"`
		Network            string `toml:"listen"`
		EndPort            int    `toml:"port"`
		MonitorEndPort     int    `toml:"monitorPort"`
		ServerReadTimeout  int64  `toml:"serverReadTimeout"`
		ServerWriteTimeout int64  `toml:"serverWriteTimeout"`
		FuncTimeThreshold  int64  `toml:"funcTimeThreshold"`
		RequestProcTimeout int64  `toml:"requestProcTimeout"`
	}

	log struct {
		FrameworkLog string `toml:"frameworklogDir"`
		LogDir       string `toml:"logDir"`
		LogFile      string `toml:"logFile"`
		LogLevel     string `toml:"logLevel"`
		MaxDays      int    `toml:"maxDays"`
		MaxLines     int64  `toml:"maxLines"`
		MaxSize      int64  `toml:"maxSize"`
		ChanLen      int64  `toml:"chanLen"`
		AnalysisFile string `toml:"analysisFile"`
	}

	stat struct {
		LogPath     string `toml:"statpath"`
		LogLevel    string `toml:"statlevel"`
		NamePrefix  string `toml:"statnameprefix"`
		Filename    string `toml:"statfilename"`
		MaxFileSize int    `toml:"statmaxfilesize"`
		MaxDays     int    `toml:"statmaxdays"`
		MaxLines    int    `toml:"statmaxlines"`
		Chanlen     int64  `toml:"statchanlen"`
		Interval    int    `toml:"statinterval"`
		Rotateperm  string `toml:"rotateperm"`
		Perm        string `toml:"perm"`
	}

	database struct {
		Mysql      mysql `toml:"mysql"`
		MysqlSlave mysql `toml:"mysqlslave"`
	}

	mysql struct {
		ServerAddr   string `toml:"serveraddr"`
		User         string `toml:"user"`
		Passwd       string `toml:"passwd"`
		Database     string `toml:"database"`
		MaxOpenConns int    `toml:"maxopenconns"`
		MaxIdleConns int    `toml:"maxidleconns"`
		IdleTimeout  int    `toml:"idletimeout"`
	}

	redisinfo struct {
		ServerList  []string `toml:"serverlist"`
		Passwd      string   `toml:"passwd"`
		MaxIdle     int      `toml:"maxIdle"`
		MaxActive   int      `toml:"maxActive"`
		IdleTimeout int      `toml:"idleTimeout"`
	}

	cfgcenter struct {
		MasterServerList []string `toml:"masterAddrList"`
	}
)

func init() {
	if configFile == "" {
		configFile = "./etc/config.toml"
	}

	if _, err := toml.DecodeFile(configFile, &Config); err != nil {
		panic(fmt.Sprintf("load server config err:%s", err.Error()))
	}
	fmt.Printf("level:%s path:%s NamePrefix:%s filename:%s interval:%d\n", Config.Server.Stat.LogLevel,
		Config.Server.Stat.LogPath,
		Config.Server.Stat.NamePrefix,
		Config.Server.Stat.Filename,
		Config.Server.Stat.Interval)

}
