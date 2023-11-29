package knet

import (
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/fatih/color"
)

type SendRequest struct {
	Name     string
	Method   string
	FetchURL string
	Data     io.Reader
	From     url.Values
	Cookies  []*http.Cookie
	Headers  map[string]string
}

func (s *SendRequest) Send(callBack func(resp []byte, cookies []*http.Cookie, err error)) {
	Method := strings.ToUpper(s.Method)

	if s.Headers == nil {
		s.Headers = map[string]string{
			"User-Agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1 Edg/117.0.0.0",
		}
	}

	switch Method {
	case "POST":
		callBack(WebPost(s.FetchURL, s.Headers, s.Cookies, s.From, s.Data))
	case "GET":
		fallthrough
	default:
		callBack(WebGet(s.FetchURL, s.Headers, s.Cookies))
	}
}

func (S *SendRequest) SaveFile(fileName, method string, body []byte) error {
	// 判断存储模式
	switch method {
	case "html":
		fileName += ".html"
	case "jpg":
		fileName += ".jpg"
	default:
		fileName += "." + method
	}

	// 写入文件
	if err := os.WriteFile(fileName, body, 0644); err != nil {
		return err
	} else {
		color.Blue("File " + fileName + " is Saved!")
	}
	return nil
}

func WebGet(urlStr string, headers map[string]string, cookies []*http.Cookie) ([]byte, []*http.Cookie, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return nil, nil, err
	}

	// 设置 headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 设置 cookies
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return body, resp.Cookies(), nil
}

func WebPost(urlStr string, headers map[string]string, cookies []*http.Cookie, from_ url.Values, data io.Reader) ([]byte, []*http.Cookie, error) {
	client := &http.Client{}
	var req *http.Request
	var err error

	if from_ == nil {
		if req, err = http.NewRequest("POST", urlStr, data); err != nil {
			return nil, nil, err
		}
	} else {
		if req, err = http.NewRequest("POST", urlStr, nil); err != nil {
			return nil, nil, err
		}
		// 设置 POST 数据
		req.PostForm = from_
	}

	// 设置 headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// 设置 cookies
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	return body, resp.Cookies(), nil
}
