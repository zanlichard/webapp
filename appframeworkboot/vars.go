package appframeworkboot

import (
	"webapp/appframework"

	"gitee.com/cristiane/go-common/log"
)

// SetupVars 加载变量
func SetupVars() error {
	var err error
	appframework.ErrorLogger, err = log.GetErrLogger("err")
	if err != nil {
		return err
	}

	appframework.AccessLogger, err = log.GetAccessLogger("access")
	if err != nil {
		return err
	}

	appframework.BusinessLogger, err = log.GetBusinessLogger("business")
	if err != nil {
		return err
	}
	return nil
}
