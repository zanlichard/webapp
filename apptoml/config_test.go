package apptoml

import (
	"testing"
)

func TestInit(t *testing.T) {

	Init("..\\etc\\config.toml")
	t.Logf("level:%s path:%s NamePrefix:%s filename:%s dir:%s\n",
		Config.Server.Stat.LogLevel,
		Config.Server.Stat.LogPath,
		Config.Server.Stat.NamePrefix,
		Config.Server.Stat.Filename,
		Config.Server.Log.LogDir)

	t.Logf("dep service:%+v", Config.ConfigMng.DepServiceList)
	for _, service := range Config.ConfigMng.DepServiceList {
		t.Logf("service :%s", service)
	}

}
