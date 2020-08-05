package handles

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"schoolserver/common/ecode"
	"schoolserver/dao/db"
	"schoolserver/dao/redisDao"
	"schoolserver/http/middleware"
	"strings"
	"time"
)

var privateKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICWwIBAAKBgQCCRmoOhpORPADC4ArCYPGdZRB8u44eQvcEsodcGcFSMDRWka1z
BkNJwvjCGRG9fwywFBZ7rp7Nt6NFhOd54wxQ4yQVWDETvrSeGifNT4OuBOS/K0wj
JOlSLu4pBwNfzK9AQ51X8h8VVShMeIgFOGT7kse9gZFeM1nDrLiLvCgygwIDAQAB
AoGAMkbTknSMie3yy0Kk1FyOkrRY1xKZKAOnCniN9dT4v2vcHxbMrkoZ+OMWlu4O
4yefsWqnPKdpITKAgAlpPiiKbTMi/TdGE5lyKa8RwCwKBgIWZZmVbSUJhbu9G+Se
gdK3/MNTUmeCFG7iwvWszTL/3DLHOmEEYJpYlyTIH/AE3XUCQQCUZq/ZF/4AIEvR
gF1UpJWm007njOJpx159HUw5Z/+0y/D4MTZ45awMZoe+bXm0fIdQY6Jokrl78Qw0
2Qhyqz3NAkEA4LtHCbmnqXVDFXrihd5+ErBKLBpRxq39yH1jdp+EAzOLLkqoCBCw
DLA9TQfq4LJtn40AIE9cnkFF4pzPNm5hjwJAEBq0qpJ39fuLPsj3V+AkfV4hCe+4
AlKoZltvkis/DJe1Jrnwd141NYNK59dphbSd2pN1ZHPHvTODZ5jF2evLYQJAV5iH
6wGDmajMWi4I55c+2vf+IUys/V1KY4CEaXNp2HmZ0ZRmBKbEiF2Vt1XTtnu2AQ/L
scxdVI4quFbY6eWCfwJAIMeUwi9aXjBht4w1QsRamaMNf2O4+1lspfR5pRTbJUB/
MEF3VlQTUCFQjTDRQ5SYo5ok2HlZOVR5mSvI1vJr/w==
-----END RSA PRIVATE KEY-----`)

//退出登录
func SignOutHandler(c echo.Context) error {
	req := c.Request()
	mySchoolCookie, err := req.Cookie(middleware.SessionName)
	if err != nil {
		return FailWithMsg(c, 4001, "退出失败")
	}
	sessionID := mySchoolCookie.Value
	err = redisDao.DeleteRedisBuyKey(sessionID)
	if err != nil {
		return FailWithMsg(c, 4002, "退出失败")
	}

	return Success(c, ecode.OK, "退出成功")
}

//登录
func SigninHandle(c echo.Context) error {
	uid := c.FormValue("uname")
	pwd := c.FormValue("pwd")

	decodeBytesId, _ := base64.StdEncoding.DecodeString(uid)
	id, err := RsaDecrypt(decodeBytesId) //RSA解密
	uid = string(id)

	decodeBytesPwd, _ := base64.StdEncoding.DecodeString(pwd)
	pwdbaty, err := RsaDecrypt(decodeBytesPwd) //RSA解密
	pwd = string(pwdbaty)

	user := db.TUser{}
	user.UserId = uid
	user.Pwd = pwd
	falg, err := user.LoginByPwd()
	if err != nil {
		return FailWithMsg(c, 4001, fmt.Sprintf("登录时发生异常：%v", err))
	}
	if !falg {
		return FailWithMsg(c, 4002, "用户名或密码错误！")
	}
	//密码验证通过，将用户信息记录到 session中
	err = addUserIntoSession(c, user)
	if err != nil {
		return FailWithMsg(c, 4001, fmt.Sprintf("登录时发生异常：%v", err))
	}
	go func() {
		addLoginHistory(c, user)
	}()
	return Success(c, ecode.OK, user)
}

//用私钥进行解密
func RsaDecrypt(cipherText []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey) //将密钥解析成私钥实例
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes) //解析pem.Decode（）返回的Block指针实例
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipherText) //RSA算法解密
}

//将用户信息保存到session中
func addUserIntoSession(c echo.Context, u db.TUser) error {

	resp := c.Response()
	//生成session值
	sessionValue := middleware.SessionName + db.NowTimeNumStr() + u.UserId

	//注意 cookie 中不能包含等号，只能是字符或数字
	cookie := http.Cookie{
		Name:   middleware.SessionName,
		Value:  sessionValue,
		MaxAge: 3600,
		Path:   "/",
		Domain: "172.29.33.90:3334",
	}
	http.SetCookie(resp.Writer, &cookie)
	resp.Header().Set("Myschool-UserId", u.UserId+":忍者神龟的化身")

	buf, err := json.Marshal(u)
	if err != nil {
		return err
	}
	//向redis 中存储session信息
	err = redisDao.SaveRedisString(sessionValue, string(buf), 10*time.Minute)
	if err != nil {
		log.Info(err)
	}

	return nil
}

//记录人员登录信息到数据库
func addLoginHistory(c echo.Context, curUser db.TUser) {

	loginPwdNum := redisDao.Client.Incr("loginPwdNum").Val()
	req := c.Request()
	agentString := req.Header.Get("User-Agent")
	hostString := req.Host
	log.Info(agentString)
	log.Info(hostString)

	ug := db.TUserAgent{}
	ug.UserId = curUser.UserId
	ug.Name = curUser.Name
	ug.LoginTime = db.NowTimeStr()
	ug.UserAgent = agentString
	ug.LoginAddress = hostString
	ug.PwdLevel = (int32(loginPwdNum)) % 5

	if strings.Index(ug.UserAgent, "Windows") != -1 {
		ug.LoginEquipment += "Windows系统--"
	} else {
		ug.LoginEquipment += "其他系统--"
	}
	if strings.Index(ug.UserAgent, "MetaSr") != -1 {
		ug.LoginEquipment += "搜狗浏览器"
	} else if strings.Index(ug.UserAgent, "Chrome") != -1 {
		ug.LoginEquipment += "谷歌浏览器"
	} else if strings.Index(ug.UserAgent, "Firefox") != -1 {
		ug.LoginEquipment += "火狐浏览器"
	} else {
		ug.LoginEquipment += "其他浏览器"
	}
	db.InsertLoginHistory(&ug)
}

//获取登录历史列表
func GetLoginHistory(c echo.Context) error {
	curUser := c.Get(middleware.CtxUser).(*db.TUser)
	d := make([]db.TUserAgent, 0, 4)
	d, err := db.GetloginHistory(curUser.UserId)
	if err != nil {
		return FailWithMsg(c, 4001, fmt.Sprintf("获取用户列表时发生异常：%v", err))
	}

	log.Info(curUser)
	return Success(c, ecode.OK, d)
}
