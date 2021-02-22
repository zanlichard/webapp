package app

import (
	"context"
	"webapp/appframework"
	"webapp/appframework/code"

	"github.com/gin-gonic/gin"
	"github.com/zanlichard/beegoe/validation"
)

func MarkErrors(ctx context.Context, errors []*validation.Error) {
	for _, err := range errors {
		appframework.BusinessLogger.Error(ctx, err.Key, err.Message)
	}
	return
}

func JsonResponse(ctx *gin.Context, httpCode, retCode int, data interface{}) {
	ctx.JSON(httpCode, gin.H{
		"code": retCode,
		"msg":  code.GetMsg(retCode),
		"data": data,
	})
}

func ProtoBufResponse(ctx *gin.Context, httpCode int, data interface{}) {
	ctx.ProtoBuf(httpCode, data)
}
