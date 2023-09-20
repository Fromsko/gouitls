package knet

import (
	"net/http"
	"testing"
)

func TestSendRequest_Send(t *testing.T) {
	// 创建一个 SendRequest 实例，用于测试 Send 方法
	req := SendRequest{
		Method:   "GET",
		FetchURL: "https://example.com",
	}

	// 定义一个回调函数，处理 Send 方法的响应
	callback := func(body []byte, cookies []*http.Cookie, err error) {
		if err != nil {
			t.Errorf("SendRequest.Send() returned an error: %v", err)
		}
	}

	// 调用 Send 方法并传入回调函数
	req.Send(callback)
}

func TestSendRequest_SaveFile(t *testing.T) {
	// 创建一个 SendRequest 实例，用于测试 SaveFile 方法
	req := SendRequest{}

	// 模拟保存文件
	fileName := "testfile.html"
	body := []byte("<html><body>Hello, World!</body></html>")

	err := req.SaveFile(fileName, "html", body)
	if err != nil {
		t.Errorf("SendRequest.SaveFile() returned an error: %v", err)
	}
}
