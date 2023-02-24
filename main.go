package main

import (
	"os"

	"pay/internal"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/yaml.v3"
)

func main() {

	content, err := os.ReadFile("conf.yml")
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(content, &internal.Wechat.Conf)
	if err != nil {
		panic(err)
	}

	internal.Wechat.Init()

	e := echo.New()
	e.Use(middleware.Recover()) // 注册异常恢复中间件

	e.GET("/", internal.Api.Index)
	e.GET("/index.html", internal.Api.Index)
	e.GET("/result.html", internal.Api.Result)
	e.POST("/h5", internal.Api.H5)
	e.POST("/notice", internal.Api.Notice)

	err = e.Start("127.0.0.1:8033")
	if err != nil {
		panic(err)
	}
}
