package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"schoolserver/common/ecode"
	"schoolserver/common/utils/encoding"
	"schoolserver/dao/db"
	"schoolserver/dao/redisDao"
	"net/http"
	"strings"
	"time"
)

const (
	CtxUser    = "currUser"
	CtxClassId = "currClassId"
	CtxUserId  = "currUserId"
	CtxLeague  = "currClass"
	baseReq    = "userId,leagueId,queryLeagueId,queryUserId"
)

var SessionName = "MySchoolSESSIONId"

//我的系统 里校验是否登录
func checkLogin(c echo.Context) ecode.ECode {
	req := c.Request()
	mySchoolCookie, err := req.Cookie(SessionName)
	if err != nil {
		return ecode.LoginExpire
	}
	sessionID := mySchoolCookie.Value

	//通过从cookie中获取的 sessionid 从redis中获取 当前登录人的用户信息
	userjsonString := ""
	err = redisDao.GetRedisString(sessionID, &userjsonString)
	if err != nil {
		return ecode.LoginExpire
	}
	currUser := db.TUser{}
	err = json.Unmarshal([]byte(userjsonString), &currUser)
	if err != nil {
		return ecode.LoginExpire
	}
	//重新设置有效期
	//向redis 中存储session信息
	redisDao.SaveRedisExpire(sessionID, 30*time.Minute)

	c.Set(CtxUserId, currUser.UserId)
	c.Set(CtxUser, &currUser)
	c.Set(CtxClassId, currUser.ClassId)
	return ecode.OK
}

// IsUserLoggedIn 检查用户是否已经登录
func IsUserLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		err := checkLogin(c)
		if err != ecode.OK {
			return fail(c, err)
		}
		return next(c)
	}
}

// IsAdmin 检查用户是否登录并且是管理员权限
func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		err := checkLogin(c)
		if err != ecode.OK {
			return fail(c, err)
		}

		u := c.Get(CtxUser).(*db.TUser)

		if u.PersonType == 1 {
			return next(c)
		} else {
			return fail(c, ecode.NotAdmin)
		}
	}
}

func EncodeAESParams(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		userId := c.Get(CtxUserId).(string)
		if len(userId) == 0 {
			return ecode.ContextMissingUid
		}

		var obj map[string]string
		buf, err := ioutil.ReadAll(c.Request().Body)
		if err != nil {
			return err
		}
		err = json.Unmarshal(buf, &obj)
		if err != nil {
			return err
		}
		for key, value := range obj {
			if strings.Contains(baseReq, key) {
				s, err := encoding.DecryptByUid(value, userId)
				if err != nil {
					return ecode.DecryptParamsError
					//return echo.NewHTTPError(9091, "操作出错")
				}
				obj[key] = s
				buf, err = json.Marshal(&obj)
				if err != nil {
					return err
				}
			}
		}
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		return next(c)

	}
}

type APIMsg struct {
	RetCode int    `json:"retCode"`
	RetMsg  string `json:"retMsg"`
}

func fail(c echo.Context, e ecode.ECode) error {
	result := APIMsg{
		RetCode: e.Code,
		RetMsg:  e.Msg,
	}
	return c.JSON(http.StatusOK, result)
}
