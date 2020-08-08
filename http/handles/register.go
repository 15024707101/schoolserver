package handles

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"schoolserver/common/ecode"
	"schoolserver/dao/db"
	"schoolserver/dao/redisDao"
)

func Register(c echo.Context) error {
	u := new(db.TUser)
	if err := c.Bind(u); err != nil {
		return FailWithMsg(c, 4001, fmt.Sprintf("注册时发生异常：%v", err))
	}
	//设置默认值
	if len(u.Mobile)==0{
		u.Mobile="15555555555"
	}
	if len(u.Pwd)==0{
		u.Pwd="123"
	}
	if len(u.ClassId)==0{
		u.ClassId="1"
	}
	if u.Age==0{
		userAge := redisDao.Client.Incr("UserAge").Val()
		u.Age=int32(userAge)
	}
	if u.PersonType==0{
		u.PersonType=1
	}
	u.CreateTime=db.NowTimeStr()
	u.Status=1
	db.InsertUser(u)

	return Success(c, ecode.OK, u)
}
