package routers

import (
	"github.com/labstack/echo/v4"
	m "github.com/labstack/echo/v4/middleware"
	"schoolserver/models"
)

func jwtMiddleware() echo.MiddlewareFunc {
	if len(models.JWTSigningKey) <= 0 {
		panic("please set jwt key")
	}
	return m.JWTWithConfig(m.JWTConfig{
		Claims:     &models.JwtCustomClaims{},
		ContextKey: "token",
		SigningKey: []byte(models.JWTSigningKey),
	})
}
