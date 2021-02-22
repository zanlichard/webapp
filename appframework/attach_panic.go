package appframework

import (
	"context"
	"runtime/debug"
)

func AttachPanicHandle(f func()) func() {
	return func() {
		defer func() {
			if err := recover(); err != nil {
				ErrorLogger.Errorf(context.Background(), "goroutine panic: %v, stacktrace:%v", err, string(debug.Stack()))
			}
		}()
		f()
	}
}
