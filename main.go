package main

import (
	"flag"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	conf "schoolserver/config"
	"schoolserver/dao/db"
	"schoolserver/dao/redisDao"
	handlers "schoolserver/http/handles"
	"schoolserver/http/routers"
	"schoolserver/logger"
	"schoolserver/models"
	"time"
)

var (
	host string
	port string
	prod bool
)

func init() {

	var configFile = flag.String("conf", "config/config.yaml", "configure file path")

	flag.StringVar(&port, "port", "3334", "listen port")
	flag.StringVar(&host, "host", "127.0.0.1", "listen host")

	var configs conf.AllConfig

	flag.Parse()
	content, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic("Configure File Read Error!")
	}

	if err := yaml.Unmarshal(content, &configs); err != nil {
		panic("unmarshal yaml format config failed")
	}

	// 初始化 log
	logger.InitLog(&configs)
	logger.Info("logger init succeed")

	// 初始化 es
	/*search.InitES(&configs)
	logger.Info("search init succeed")
	err = search.InitEsConn()
	mustOk(err)*/

	// 初始化 Mysql
	err = db.InitMysqlDB(&configs)
	mustOk(err)
	logger.Info("mysql init succeed")

	// 初始化 redis
	err = db.InitRedis(&(configs.Redis))
	mustOk(err)
	logger.Info("redisDao init succeed")

	//api.BackendAddr = configs.Sys.BackendAddr

	models.SetSingKey(configs.Sys.JwtKey)


	// 载入基本定义
	/*err = dao.LoadFromApi()
	mustOk(err)
	logger.Info("load definedIds succeed")*/

	if configs.Sys.ServiceMode == "prod" {
		prod = true
	}
}

func main() {
	e := ServerNew()
	e.HideBanner = true
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	go StartServer(e, host, port, prod)
	WaitForSignals()

	// close all connection
	redisDao.Client.Close()
	db.Close()

}

func ServerNew() (e *echo.Echo) {
	e = echo.New()
	e.Static("/assets", "assets")
	e.Server.ReadHeaderTimeout = 10 * time.Second
	e.Server.ReadTimeout = 12 * time.Second
	e.Server.WriteTimeout = 60 * time.Second
	e.Server.IdleTimeout = 0
	e.Server.SetKeepAlivesEnabled(false)
	// 自定义Error提示
	e.HTTPErrorHandler = handlers.HTTPErrorHandler
	routers.RegisterRouters(e)
	RegisterMiddleware(e)
	//PrintRoutes(e)

	e.Any("/v1/maotai/pprof/cmdline", handlers.Cmdlineec)
	e.Any("/v1/maotai/pprof/profile", handlers.Profileec)
	e.Any("/v1/maotai/pprof/symbol", handlers.Symbolec)
	e.Any("/v1/maotai/pprof/trace", handlers.Traceec)

	//PrintRoutes(e)
	return e
}

// StartServer 启动一个API服务
func StartServer(e *echo.Echo, host string, port string, prod bool) {
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", host, port)))
}

func mustOk(err error) {
	if err != nil {
		panic(err)
	}
}
func RegisterMiddleware(e *echo.Echo) {
	// 设置日志
	e.HideBanner = true
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://127.0.0.1:8077",
			"http://127.0.0.1:8080",
			"http://localhost:8080",
			"*",
			"http://172.29.33.90:3333",
			"http://172.29.33.90:8080",
		},
		AllowMethods:     []string{echo.GET, echo.PUT, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowCredentials: true,
	}))

}
