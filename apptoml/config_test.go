package apptoml

import (
        "testing"
)

func TestInit(t *testing.T) {
        t.Logf("level:%s path:%s NamePrefix:%s filename:%s dir:%s\n",
        	     Config.Server.Stat.LogLevel,
        	     Config.Server.Stat.LogPath,
			     Config.Server.Stat.NamePrefix,
			     Config.Server.Stat.Filename,
			     Config.Server.Log.LogDir)

}
