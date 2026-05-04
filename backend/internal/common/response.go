package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 响应状态码
const (
	CodeSuccess      = 0
	CodeParamError   = 400
	CodeUnauthorized = 401
	CodeForbidden    = 403
	CodeNotFound     = 404
	CodeServerError  = 500
)

// Response 统一响应结构体
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  "success",
		Data: data,
	})
}

// Fail 失败响应
func Fail(c *gin.Context, code int, msg string) {
	c.JSON(http.StatusOK, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}

// ParamError 参数错误
func ParamError(c *gin.Context, msg string) {
	Fail(c, CodeParamError, msg)
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context) {
	Fail(c, CodeUnauthorized, "unauthorized")
}

// ServerError 服务器错误
func ServerError(c *gin.Context) {
	Fail(c, CodeServerError, "server error")
}
