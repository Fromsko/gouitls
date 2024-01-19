package reply

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JsonMsg 通用返回格式
type JsonMsg struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Err  string `json:"err"`
	flag bool
}

// Option 函数选项式
type Option func(JsonMsg)

// WithMsg 初始化数据
func WithMsg(msg string) Option {
	return func(jm JsonMsg) {
		jm.Msg = msg
	}
}

// WithErr 初始化错误
func WithErr(err string) Option {
	return func(jm JsonMsg) {
		jm.Err = err
	}
}

// WithCode 初始化状态码
func WithCode(code int) Option {
	return func(jm JsonMsg) {
		jm.Code = code
	}
}

// WithFlag 初始化标志位
func WithFlag(is bool) Option {
	return func(jm JsonMsg) {
		jm.flag = is
	}
}

// ReplyData 定义数据返回格式
func Client(ctx *gin.Context, jm *JsonMsg, opts ...Option) {
	if jm == nil {
		jm = &JsonMsg{}
		for _, opt := range opts {
			opt(*jm)
		}
	}

	if jm.flag {
		ctx.JSON(jm.Code, jm)
	} else {
		ctx.JSON(http.StatusOK, jm)
	}
}
