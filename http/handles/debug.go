package handles

import (
	"github.com/labstack/echo/v4"
)

func Ping(ctx echo.Context) error {

	return ctx.String(200, "hello")
}

func Login(ctx echo.Context) error {

	//token, err := generateToken(loginResp.UserId, curLeagueId)
	//if err != nil {
	//	logger.Error(err)
	//	return Fail(ctx, ecode.JWTGenError)
	//}
	token := "hello"
	return ctx.String(200, token)
}
