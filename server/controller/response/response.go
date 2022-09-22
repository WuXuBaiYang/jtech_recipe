package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// P 输出数据结构体
type P struct {
	// HTTP 状态码 & 自定义状态码
	Code int
	// 输出消息
	Message string
	// 输出自定义数据
	Data any
}

// R Response 结构体
type R struct {
	// 输出类型，例如：JSON、XML、YAML
	T string
	// 输出数据
	P P
}

// 报文基础结构
func response(code int, message string, data any) (int, gin.H) {
	response := gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	}
	return code, response
}

// 成功报文构造
func success(params P) (int, gin.H) {
	code := http.StatusOK
	message := http.StatusText(code)
	data := params.Data

	fmt.Println(params)

	if params.Code != 0 {
		code = params.Code
	}

	if params.Message != "" {
		message = params.Message
	}

	return response(code, message, data)
}

// 失败报文构造
func fail(params P) (int, gin.H) {
	code := http.StatusInternalServerError
	message := http.StatusText(code)
	data := params.Data

	fmt.Println(params)

	if params.Code != 0 {
		code = params.Code
	}

	if params.Message != "" {
		message = params.Message
	}

	return response(code, message, data)
}

// 按类型写入报文
func write(context *gin.Context, contextType string, code int, res gin.H) {
	switch contextType {
	case "IndentedJSON":
		context.IndentedJSON(code, res)
	case "SecureJSON":
		context.SecureJSON(code, res)
	case "JSON":
		context.JSON(code, res)
	case "AsciiJSON":
		context.AsciiJSON(code, res)
	case "PureJSON":
		context.PureJSON(code, res)
	case "XML":
		context.XML(code, res)
	case "YAML":
		context.YAML(code, res)
	case "ProtoBuf":
		context.ProtoBuf(code, res)
	}
}

// Fail 失败报文
func Fail(context *gin.Context, response R) {
	// 默认输出 JSON 数据
	contextType := "JSON"
	if response.T != "" {
		contextType = response.T
	}
	// 构建输出数据
	code, res := fail(response.P)
	write(context, contextType, code, res)
}

// Success 成功报文
func Success(context *gin.Context, response R) {
	// 默认输出 JSON 数据
	contextType := "JSON"
	if response.T != "" {
		contextType = response.T
	}
	// 构建输出数据
	code, res := success(response.P)
	write(context, contextType, code, res)
}
