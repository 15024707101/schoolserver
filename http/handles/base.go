package handles

import (
	"github.com/labstack/echo/v4"
	"schoolserver/common/ecode"
	"schoolserver/common/utils"
	"schoolserver/logger"
	"net/http"
)

var panicErr error

func RespCheck(i interface{}, c echo.Context) error {
	logger.Debug("ressult: %+v \n", i)

	// 状态码 c
	status, err := utils.GetField(i, "HeaderResp[0].Status")
	if err != nil {
		logger.Error(err)
	}
	logger.Debug("响应状态： %v\n", status)

	desc, err := utils.GetField(i, "HeaderResp[0].Desc")
	if err != nil {
		logger.Error(err)
	}
	logger.Debug("响应消息：%v \n", desc)

	switch status.(int32) {
	case 0:
		// 现在需要根据后端返回的描述判断
		switch desc.(string) {
		case "您的登录已失效！请重新登录":
			return Fail(c, ecode.LoginExpire)
		case "Session is expried or not existed.":
			return Fail(c, ecode.LoginExpire)
		default:
			return Fail(c, ecode.LoginExpire)
		}

	case 2:
		return Fail(c, ecode.ParamsError)
	case 3:
		return Fail(c, ecode.ParamsError)
	default:
		return Success(c, ecode.OK, i)
	}
}

type APIResult struct {
	RetCode int         `json:"retCode"`
	RetMsg  string      `json:"retMsg"`
	Data    interface{} `json:"results"` // tag中带有",string"选项，那么这个字段在输出到JSON的时候会把该字段对应的值转换成JSON字符串
}

// Error 返回错误码
type APIMsg struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
}

func Success(c echo.Context, e ecode.ECode, data interface{}) error {
	result := APIResult{
		RetCode: e.Code,
		RetMsg:  e.Msg,
		Data:    data,
	}

	return c.JSON(http.StatusOK, result)
	/*b, err := json.Marshal(result)
	if err != nil {
		return err
	}
	return c.JSONBlob(http.StatusOK, b)*/
}

func Fail(c echo.Context, e ecode.ECode) error {
	result := APIMsg{
		RetCode: e.Code,
		RetMsg:  e.Msg,
	}

	if e.Code == 404 {
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusOK, result)
}

func FailAttachMsg(c echo.Context, e ecode.ECode, msg string) error {
	e.Msg += msg
	return Fail(c, e)
}

// FailWithMsg 兼容后台返回的信息 现在可以用来将api返回的msg发送给用户
func FailWithMsg(c echo.Context, code int, msg string) error {
	result := APIMsg{
		RetCode: code,
		RetMsg:  msg,
	}

	if code == 404 {
		return c.JSON(http.StatusNotFound, result)
	}
	return c.JSON(http.StatusOK, result)
}
func zhtjDefer() {
	if err := recover(); err != nil {
		logger.Error("ZHTJ-PANIC:", err)
		e, ok := err.(error)
		if ok {
			panicErr = e
		}
		return
	}
}
