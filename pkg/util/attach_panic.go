package util

import (
	"context"
	"runtime/debug"
	"webapp/globalconfig"
)

func AttachPanicHandle(f func()) func() {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				globalconfig.ErrorLogger.Errorf(context.Background(), "goroutine panic: %v, stacktrace:%v", err, string(debug.Stack()))
			}
		}()
		f()
	}
}
