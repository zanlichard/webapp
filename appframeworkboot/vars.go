package appframeworkboot

import (
	"webapp/globalconfig"
	"gitee.com/cristiane/go-common/log"
)

// SetupVars 加载变量
func SetupVars() error {
	var err error
	globalconfig.ErrorLogger, err = log.GetErrLogger("err")
	if err != nil {
		return err
	}

	globalconfig.AccessLogger, err = log.GetAccessLogger("access")
	if err != nil {
		return err
	}
	return nil
}