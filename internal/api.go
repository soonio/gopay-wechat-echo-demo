package internal

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-pay/gopay"
	sdk "github.com/go-pay/gopay/wechat"
	"github.com/labstack/echo/v4"
)

var Api = new(api)

type api struct{}

func (h *api) Index(c echo.Context) error {
	return c.File("template/index.html")
}

func (h *api) Result(c echo.Context) error {
	return c.File("template/result.html")
}

func (h *api) H5(c echo.Context) error {
	path, err := Logic.H5("HiTest商品", "Seer"+time.Now().Format("20060102150405"), 1, c.RealIP())
	if err != nil {
		return Error(c, 899, err.Error())
	}
	return Success(c, path+"&redirect_url="+url.QueryEscape(Wechat.Conf.Uri("/result.html")))
}

func (h *api) Notice(c echo.Context) error {
	bm, err := sdk.ParseNotifyToBodyMap(c.Request())
	rsp := new(sdk.NotifyResponse)
	if err == nil {
		go Logic.Callback(bm)
		rsp.ReturnCode = gopay.SUCCESS
		rsp.ReturnMsg = gopay.OK
	} else {
		rsp.ReturnCode = gopay.FAIL
		rsp.ReturnMsg = err.Error()
	}

	return c.String(http.StatusOK, rsp.ToXmlString())
}
