package models

import "github.com/dgrijalva/jwt-go"

var JWTSigningKey = ""

func SetSingKey(key string) {
	JWTSigningKey = key
}

type SigninReqParam struct {
	Brand      string
	Model      string
	Version    string
	UserAgent  string
	System     string
	PlatForm   string
	AppId      string
	WxOpenId   string
	QrCode     string `json:"qrCode"`
	WxNickName string `json:"wxNickName"`
	//// 验证码
	//CaptchaID string `json:"captcha_Id"`
	//Captcha   string `json:"captcha"`
}

type ScanReqParam struct {
	Brand    string `json:"brand"`
	Model    string `json:"model"`
	Version  string `json:"version"`
	System   string `json:"system"`
	PlatForm string `json:"platForm"`
	AppId    string `json:"appId"`
	WxOpenId string `json:"wxOpenId"`
	Uuid     string `json:"uuid"`
}
type JwtCustomClaims struct {
	UserId   string `json:"userId"`
	LeagueId string `json:"leagueId"`
	jwt.StandardClaims
}
