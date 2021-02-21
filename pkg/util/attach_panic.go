package util

import (
	"context"
	"webapp/webconfig"
	"runtime/debug"
)

func AttachPanicHandle(f func()) func() {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				webconfig.ErrorLogger.Errorf(context.Background(), "goroutine panic: %v, stacktrace:%v", err, string(debug.Stack()))
			}
		}()
		f()
	}
}
