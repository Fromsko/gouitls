package v2

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/tidwall/gjson"
)

// Option 类型，表示可变参数的配置函数
type Option func(*Request)

// 定义一个回调函数类型，用于执行请求
type RequestFunc func() (resp *http.Response, body []byte, err error)

// IRequest 接口，定义请求的行为
type IRequest interface {
	WithURL(url string) IRequest
	WithMethod(method string) IRequest
	WithData(data io.Reader) IRequest
	WithForm(form url.Values) IRequest
	WithProxy(url string) IRequest
	WithCookies(cookies []*http.Cookie) IRequest
	WithHeaders(headers map[string]string) IRequest
	Send(callBack func(IResponse, error))
}

// IResponse 接口，定义响应的行为
type IResponse interface {
	StatusCode() int
	Json() *gjson.Result
	Text() string
	Bytes() []byte
	RaiseCode(code int) error
	WriteFile(fileName, perm, dir string) error
	WriteJson(fileName, perm, dir string) error
}

// Request 实现 IRequest 接口
type Request struct {
	fetchURL string
	method   string
	data     io.Reader
	form     url.Values
	client   *http.Client
	cookies  []*http.Cookie
	headers  map[string]string
}

// NewRequest 支持两种方式构造请求
func NewRequest(options ...Option) IRequest {
	req := &Request{
		headers: make(map[string]string),
	}

	// 应用所有 Option 函数
	for _, opt := range options {
		opt(req)
	}

	return req
}

// WithProxy 配置方法为链式调用
func (r *Request) WithProxy(uri string) IRequest {
	if r.client == nil {
		r.client = &http.Client{}
	}

	proxyParsed, _ := url.Parse(uri)
	r.client.Transport = &http.Transport{
		Proxy: http.ProxyURL(proxyParsed),
	}

	return r
}

// WithMethod 配置方法为链式调用
func (r *Request) WithMethod(method string) IRequest {
	r.method = method
	return r
}

// WithURL 配置 URL 为链式调用
func (r *Request) WithURL(url string) IRequest {
	r.fetchURL = url
	return r
}

// WithData 配置数据为链式调用
func (r *Request) WithData(data io.Reader) IRequest {
	r.data = data
	return r
}

// WithForm 配置表单数据为链式调用
func (r *Request) WithForm(form url.Values) IRequest {
	r.form = form
	return r
}

// WithHeaders 配置 headers 为链式调用
func (r *Request) WithHeaders(headers map[string]string) IRequest {
	r.headers = headers
	return r
}

// WithCookies 配置 cookies 为链式调用
func (r *Request) WithCookies(cookies []*http.Cookie) IRequest {
	r.cookies = cookies
	return r
}

func (r *Request) send() (resp *http.Response, body []byte, err error) {

	var req *http.Request

	if r.headers == nil {
		r.headers = map[string]string{
			"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1 Edg/117.0.0.0",
		}
	}

	if r.data == nil && r.form == nil {
		req, err = http.NewRequest(r.method, r.fetchURL, nil)
	} else if r.data != nil {
		req, err = http.NewRequest(r.method, r.fetchURL, r.data)
	} else if r.form != nil {
		req, err = http.NewRequest(r.method, r.fetchURL, nil)
		req.PostForm = r.form
	}

	for key, value := range r.headers {
		req.Header.Set(key, value)
	}

	for _, cookie := range r.cookies {
		req.AddCookie(cookie)
	}

	resp, err = r.client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return resp, body, nil
}

// Send 发送请求
func (r *Request) Send(callBack func(resp IResponse, err error)) {
	resp, body, err := r.send()

	callBack(
		&Response{
			body: body,
			resp: resp,
		},
		err,
	)
}

// Response 实现 IResponse 接口
type Response struct {
	resp *http.Response
	body []byte
}

func (r *Response) StatusCode() int {
	return r.resp.StatusCode
}

func (r *Response) Json() *gjson.Result {
	jsonResult := gjson.ParseBytes(r.body)
	return &jsonResult
}

func (r *Response) Text() string {
	return string(r.body)
}

func (r *Response) Bytes() []byte {
	return r.body
}

func (r *Response) RaiseCode(code int) error {
	if r.resp.StatusCode == code {
		return fmt.Errorf("Error: StatusCode %d", code)
	}
	return nil
}

func (r *Response) WriteFile(fileName, perm, dir string) error {
	if dir != "" {
		fileName = dir + "/" + fileName
	}
	err := os.WriteFile(fileName, r.body, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (r *Response) WriteJson(fileName, perm, dir string) error {
	var jsonBody bytes.Buffer
	if err := json.Indent(&jsonBody, r.body, "", "  "); err != nil {
		return err
	}
	return r.WriteFile(fileName, perm, dir)
}

// Option 函数实现，支持可变参数调用
func WithMethod(method string) Option {
	return func(r *Request) {
		r.method = method
	}
}

func WithURL(url string) Option {
	return func(r *Request) {
		r.fetchURL = url
	}
}

func WithData(data io.Reader) Option {
	return func(r *Request) {
		r.data = data
	}
}

func WithForm(form url.Values) Option {
	return func(r *Request) {
		r.form = form
	}
}

func WithProxy(uri string) Option {
	return func(r *Request) {
		if r.client == nil {
			r.client = &http.Client{}
		}

		proxyParsed, _ := url.Parse(uri)
		r.client.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyParsed),
		}
	}
}
