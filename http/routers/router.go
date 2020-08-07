package routers

import (
	"github.com/labstack/echo/v4"
	"schoolserver/http/handles"
	"schoolserver/http/middleware"
)

func RegisterRouters(router *echo.Echo) {
	unlogin := router.Group("/unlogin")
	{

		unlogin.POST("/signin", handles.SigninHandle)
		unlogin.POST("/decode", handles.GetDecodeString)
		unlogin.POST("/uploadImg", handles.UploadImg)
	}
	// testing
	test := router.Group("/test")
	{
		test.Use(jwtMiddleware(), middleware.IsUserLoggedIn)
		test.GET("/hello", handles.Ping)
	}

	admin := router.Group("/center")
	{
		admin.Use(middleware.IsAdmin)
		admin.POST("/userList", handles.GetUserList)
		admin.POST("/classList", handles.GetClassList)
		admin.POST("/loginHistoryList", handles.GetLoginHistory)
	}
	router.POST("signout", handles.SignOutHandler, middleware.IsUserLoggedIn)

	sign := router.Group("/api/v1")
	{
		sign.Use(jwtMiddleware(), middleware.IsUserLoggedIn, middleware.EncodeAESParams)
	}

}
