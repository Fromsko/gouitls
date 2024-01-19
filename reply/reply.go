package reply

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// JsonMsg 通用返回格式
type JsonMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Err  *string     `json:"err,omitempty"`
	Data interface{} `json:"data,omitempty"`
	flag bool
}

// Option 函数选项式
type Option func(*JsonMsg)

// WithMsg 初始化数据
func WithMsg(msg string) Option {
	return func(jm *JsonMsg) {
		jm.Msg = msg
	}
}

// WithErr 初始化错误
func WithErr(err string) Option {
	return func(jm *JsonMsg) {
		jm.Err = &err
	}
}

// WithCode 初始化状态码
func WithCode(code int) Option {
	return func(jm *JsonMsg) {
		jm.Code = code
	}
}

// WithFlag 初始化标志位
func WithFlag(is bool) Option {
	return func(jm *JsonMsg) {
		jm.flag = is
	}
}

// 增加其他数据
func WithOther(other any) Option {
	return func(jm *JsonMsg) {
		jm.Data = other
	}
}

// ReplyData 定义数据返回格式
func Client(ctx *gin.Context, jm *JsonMsg, opts ...Option) {
	// 如果 jm 为 nil，则创建一个新的 JsonMsg
	if jm == nil {
		jm = &JsonMsg{}
	}

	// 应用选项
	for _, opt := range opts {
		opt(jm)
	}

	// 设置默认的 HTTP 状态码
	statusCode := http.StatusOK
	if !jm.flag {
		statusCode = jm.Code
	}

	// 返回 JSON 响应
	ctx.JSON(statusCode, jm)
}

// Succeed 正确数据
func Succeed(ctx *gin.Context, msg string, others ...gin.H) {
	rt := gin.H{
		"code": http.StatusOK,
		"msg":  msg,
	}

	if len(others) != 0 {
		for _, other := range others {
			for k, v := range other {
				rt[k] = v
			}
		}
	}

	ctx.JSON(http.StatusOK, rt)
}

// Failed 错误数据
func Failed(ctx *gin.Context, msg, err string, others ...gin.H) {
	rt := gin.H{
		"code": http.StatusBadRequest,
		"msg":  msg,
		"err":  err,
	}

	if len(others) != 0 {
		for _, other := range others {
			for k, v := range other {
				rt[k] = v
			}
		}
	}

	ctx.JSON(rt["code"].(int), rt)
}
