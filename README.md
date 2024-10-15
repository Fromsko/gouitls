# gouitls
一些开箱即用的小工具包

## 安装

```shell
go get github.com/Fromsko/gouitls
```

## 使用

[`knet/kent.go`](./knet/kent.go)
> 简单的 `net/http` 封装

+ 每日一句
	```go
	package api

	import (
		"net/http"

		"github.com/Fromsko/gouitls/knet"
		"github.com/fatih/color"
		"github.com/tidwall/gjson"
	)

	// GetEveryDay 获取每日一句
	func GetEveryDay() string {
		var equiangular string
		Spider := knet.SendRequest{
			FetchURL: "http://open.iciba.com/dsapi/?date",
		}
		Spider.Send(func(resp []byte, cookies []*http.Cookie, err error) {
			if err != nil {
				color.Red("获取每日一句失败")
				equiangular = "千里之堤, 始于足下。"
				return
			}
			equiangular = gjson.Get(string(resp), "note").String()
		})
		return equiangular
	}
	```

+ 和天气 数据
	```go
	package api

	import (
		"fmt"
		"net/http"
		"net/url"
		"time"

		"notify/enum"

		"github.com/fatih/color"
		"github.com/tidwall/gjson"

		"github.com/Fromsko/gouitls/knet"
	)

	var (
		// 和天气 Key
		WeatherKey = ""
		// 和天气 API 地址
		WeatherUrl = ""
		// 和天气 城市 API 地址
		WeatherCityList = ""
	)

	// WeatherObject 天气对象
	type WeatherObject struct {
		Local       string
		WeatherID   string
		WeatherInfo struct {
			Text string
			Temp string
		}
		WeatherStatus int64
		WeatherDate   time.Time
	}

	// GetWeatherID 获取 天气ID
	func (w *WeatherObject) GetWeatherID() {
		weather := knet.SendRequest{
			FetchURL: fmt.Sprintf(
				WeatherCityList,
				url.QueryEscape(w.Local),
				WeatherKey,
			),
		}
		weather.Send(func(resp []byte, cookies []*http.Cookie, err error) {
			statusCode := gjson.GetBytes(resp, "code").Int()
			if statusCode != 200 || err != nil {
				color.Red("天气请求失败!")
			} else {
				location := gjson.GetBytes(resp, "location").Array()[0]
				ID := location.Get("id").String()
				w.WeatherID = ID
			}
			w.WeatherStatus = statusCode
		})
	}

	// GetWeatherInfo 获取天气信息
	func (w *WeatherObject) GetWeatherInfo() {
		weather := knet.SendRequest{
			FetchURL: fmt.Sprintf(
				WeatherUrl,
				w.WeatherID,
				WeatherKey,
			),
		}
		weather.Send(func(resp []byte, cookies []*http.Cookie, err error) {
			statusCode := gjson.GetBytes(resp, "code").Int()
			if statusCode != 200 || err != nil {
				color.Red("天气请求失败!")
				w.WeatherInfo.Text = "未知"
				w.WeatherInfo.Temp = "未知"
			} else {
				w.WeatherInfo.Text = gjson.GetBytes(resp, "now.text").String()
				w.WeatherInfo.Temp = gjson.GetBytes(resp, "now.temp").String()
			}
		})
	}

	// SearchWeather 获取指定地方的天气
	func SearchWeather(local string) *WeatherObject {
		weather := &WeatherObject{
			Local: local,
		}
		weather.GetWeatherID()
		weather.GetWeatherInfo()
		return weather
	}

	```

## V2 版本
> 推荐使用`V2`版本, 支持函数调用和链式调用

`安装`
```shell
go get github.com/Fromsko/gouitls
```

`使用`
```go
package main

import (
	"fmt"
	"net/http"

	. "github.com/Fromsko/gouitls/knet/v2"
)

func main() {
	spider := NewRequest(
		WithMethod(http.MethodGet),
		WithProxy("http://127.0.0.1:7890"),
		WithURL("https://uapis.cn/api/hotlist?type=history"),
	)
	spider.Send(func(resp IResponse, err error) {
		fmt.Printf("resp.Json(): %v\n", resp.Json())
		fmt.Printf("resp.StatusCode(): %v\n", resp.StatusCode())
		resp.WriteJson(
			"history.json",
			"0644",
			"",
		)
	})
}
```
