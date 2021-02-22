package appframeworkboot

import (
	"gitee.com/cristiane/go-common/log"
	"webapp/globalconfig"
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

	globalconfig.BusinessLogger, err = log.GetBusinessLogger("business")
	if err != nil {
		return err
	}
	return nil
}