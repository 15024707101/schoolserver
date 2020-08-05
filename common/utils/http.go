package utils

import (
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

// FetchRealUrl 获取链接真实的URL（获取重定向一次的结果URL）
func FetchRealUrl(uri string) (realUrl string) {

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			realUrl = req.URL.String()
			return errors.New("util fetch real url")
		},
	}

	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return uri
}

const XRequestedWith = "X-Requested-With"

func IsAjax(ctx echo.Context) bool {
	if ctx.Request().Header.Get(XRequestedWith) == "XMLHttpRequest" {
		return true
	}
	return false
}
