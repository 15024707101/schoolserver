package handles

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"schoolserver/common/ecode"
	"schoolserver/common/utils/encoding"
	"schoolserver/dao/db"
	"schoolserver/http/middleware"
	"strings"
)

func GetUserList(c echo.Context) error {
	curUser := c.Get(middleware.CtxUser).(*db.TUser)
	d := make([]db.TUser, 0, 4)
	d, err := db.GetUserList()
	if err != nil {
		return FailWithMsg(c, 4001, fmt.Sprintf("获取用户列表时发生异常：%v", err))
	}

	log.Info(curUser)
	return Success(c, ecode.OK, d)
}

func GetClassList(c echo.Context) error {
	curUser := c.Get(middleware.CtxUser).(*db.TUser)
	d := make([]db.TClass, 0, 4)
	d, err := db.GetClassList()
	if err != nil {
		return FailWithMsg(c, 4001, fmt.Sprintf("获取用户列表时发生异常：%v", err))
	}

	log.Info(curUser)
	return Success(c, ecode.OK, d)
}

/*解密功能*/
func GetDecodeString(c echo.Context) error {
	str := c.FormValue("str")
	uId := c.FormValue("key")
	ddtype := c.FormValue("type")

	var uu string

	rowDatas := []string{str}
	switch ddtype {
	case "1", "a":
		var sss strings.Builder
		var e error
		for i, v := range rowDatas {
			ss, err := encoding.DecryptRand(v)
			if err != nil {
				e = err
				ss = "解密失败"
			}
			sss.WriteString(ss)
			if i == 0 {
				ul := strings.Split(ss, "_")
				uu = ul[0]
			}
			if i < len(rowDatas)-1 {
				sss.WriteString(",")
			}
		}
		con := sss.String()
		if e != nil {
			fmt.Println("其中 解密错误信息：\n", e)
			return FailWithMsg(c, 4001, "解密错误！")
		}
		if len(con) <= 0 {
			fmt.Println("解密错误", uu)
			return FailWithMsg(c, 4002, "解密错误！")
		} else {
			fmt.Println("解密后的内容为，默认记住当前的UID，方便后续使用 [userId_leagueId]：\n", con)
			return Success(c, ecode.OK, con)
		}

	case "2", "b":
		if len(uId) <= 0 {
			var sss strings.Builder
			var e error
			for i, v := range rowDatas {
				ss, err := encoding.DecryptByUid(v, "")
				if err != nil {
					e = err
					ss = "解密失败"
				}
				sss.WriteString(ss)
				if i < len(rowDatas)-1 {
					sss.WriteString(",")
				}
			}
			con := sss.String()
			if e != nil {
				fmt.Println("其中 解密错误信息：\n", e)
				return FailWithMsg(c, 4004, "解密错误！")
			}
			if len(con) <= 0 {
				fmt.Println("解密错误")
				return FailWithMsg(c, 4003, "解密错误！")
			} else {
				fmt.Println("解密后的内容为：\n", con)
				return Success(c, ecode.OK, con)
			}

		} else {
			var sss strings.Builder
			var e error
			for i, v := range rowDatas {
				ss, err := encoding.DecryptByUid(v, uId)
				if err != nil {
					e = err
					ss = "解密失败"
				}
				sss.WriteString(ss)
				if i < len(rowDatas)-1 {
					sss.WriteString(",")
				}
			}
			con := sss.String()
			if e != nil {
				fmt.Println("其中 解密错误信息：\n", e)
				return FailWithMsg(c, 4005, "解密错误！")
			}
			if len(con) < 0 {
				fmt.Println("解密错误")
				return FailWithMsg(c, 4006, "解密错误！")
			} else {
				fmt.Println("解密后的内容为：\n", con)
				return Success(c, ecode.OK, con)
			}
		}
	case "3", "c":
		var sss strings.Builder
		var e error
		for i, v := range rowDatas {
			ss, err := encoding.DecryptByUid(v, "")
			if err != nil {
				e = err
				ss = "解密失败"
			}
			sss.WriteString(ss)
			if i < len(rowDatas)-1 {
				sss.WriteString(",")
			}
		}
		con := sss.String()
		if e != nil {
			fmt.Println("其中 解密错误信息：\n", e)
			return FailWithMsg(c, 4007, "解密错误！")
		}
		if len(con) < 0 {
			fmt.Println("解密错误")
			return FailWithMsg(c, 4008, "解密错误！")
		} else {
			fmt.Println("解密后的内容为：\n", con)
			return Success(c, ecode.OK, con)
		}

	case "4", "d":
		var sss strings.Builder
		var e error
		for i, v := range rowDatas {
			ss, err := encoding.DecryptExpireData(v)
			if err != nil {
				e = err
				ss = "解密失败"
			}
			sss.WriteString(ss)
			if i < len(rowDatas)-1 {
				sss.WriteString(",")
			}
		}
		con := sss.String()
		if e != nil {
			fmt.Println("其中 解密错误信息：\n", e)
			return FailWithMsg(c, 4009, "解密错误！")
		}
		if len(con) < 0 {
			fmt.Println("解密错误")
			return FailWithMsg(c, 4010, "解密错误！")
		} else {
			fmt.Println("解密后的内容为：\n", con)
			return Success(c, ecode.OK, con)
		}

	default:
		fmt.Println("警告：类型不在范围内")
	}
	return FailWithMsg(c, 4002, "解密错误！")
}
