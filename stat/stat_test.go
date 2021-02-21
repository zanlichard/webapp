package stat

import (
	"testing"
	"time"
)

func TestStat(t *testing.T) {
	logconfig := new(LoggerParam)
	logconfig.Level = "info"
	logconfig.Path = "."
	logconfig.NamePrefix = "test"
	logconfig.Filename = "stat.log"
	logconfig.Maxfilesize = 10000
	logconfig.Maxdays = 7
	logconfig.Maxlines = 10000
	logconfig.Chanlen = 10000
	Init(*logconfig, 20)
	SetDelayUp(20, 50, 100)
	StatProc()
	time.Sleep(21 * time.Second)
	//stat.SendStatItem(elem)
	Exit()

}
