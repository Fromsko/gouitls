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

// WithMsg 返回数据
func WithMsg(msg string) Option {
	return func(jm JsonMsg) {
		jm.Msg = msg
	}
}

// WithErr 返回错误
func WithErr(err string) Option {
	return func(jm JsonMsg) {
		jm.Err = err
	}
}

// WithCode 返回状态码
func WithCode(code int) Option {
	return func(jm JsonMsg) {
		jm.Code = code
	}
}

// ReplyData 定义数据返回格式
func Client(ctx *gin.Context, opts ...Option) {
	c := &JsonMsg{}

	for _, opt := range opts {
		opt(*c)
	}

	if c.flag {
		ctx.JSON(c.Code, c)
	} else {
		ctx.JSON(http.StatusOK, c)
	}
}
