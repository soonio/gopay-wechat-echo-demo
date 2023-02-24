package internal

import (
	"fmt"
	sdk "github.com/go-pay/gopay/wechat"
	"github.com/labstack/echo/v4"
	"net/http"
)

var Wechat = new(wechat)

type wechatConfig struct {
	Host  string `json:"host"   yaml:"host"`
	AppID string `json:"appid"  yaml:"appid"`
	MchID string `json:"mch-id" yaml:"mch-id"`
	Key   string `json:"key"    yaml:"key"`
}

func (h *wechatConfig) Uri(path string) string {
	return fmt.Sprintf("%s%s", h.Host, path)
}

type wechat struct {
	Client *sdk.Client
	Conf   *wechatConfig
}

func (w *wechat) Init() {
	w.Client = sdk.NewClient(w.Conf.AppID, w.Conf.MchID, w.Conf.Key, true)
}

type Response struct {
	Code int    `json:"code"`           // 业务状态码
	Msg  string `json:"msg"`            // 业务消息
	Data any    `json:"data,omitempty"` // 数据
}

// Success 成功时响应消息
func Success(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, &Response{0, "success", data})
}

// Error 失败时响应消息
func Error(c echo.Context, code int, msg string) error {
	return c.JSON(http.StatusInternalServerError, &Response{Code: code, Msg: msg})
}
