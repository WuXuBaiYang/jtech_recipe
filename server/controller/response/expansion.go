package response

import (
	"github.com/gin-gonic/gin"
	"server/tool"
)

// FailDef 默认失败报文
func FailDef(context *gin.Context, code int, message string) {
	Fail(context, R{
		P: P{
			Code:    code,
			Message: message,
		},
	})
}

// SuccessDef 默认成功报文
func SuccessDef(context *gin.Context, data any) {
	Success(context, R{
		P: P{
			Code:    0,
			Message: "",
			Data:    data,
		},
	})
}

// FailServer 服务器异常报文
func FailServer(context *gin.Context) {
	FailDef(context, -1, "服务器异常")
}

// FailParams 请求参数异常
func FailParams(context *gin.Context, message string) {
	FailDef(context, -1, message)
}

// FailParamsDef 请求参数异常默认值
func FailParamsDef(context *gin.Context, err error) {
	msg := tool.If(gin.IsDebugging(), err.Error(), "参数异常")
	FailParams(context, msg)
}

// FailAuth 授权失效
func FailAuth(context *gin.Context, message string) {
	FailDef(context, 401, message)
}

// FailAuthDef 授权失效默认
func FailAuthDef(context *gin.Context) {
	FailAuth(context, "权限不足")
}
