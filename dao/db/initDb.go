package db

import (
	"errors"
	"fmt"
	mysql_ "github.com/go-sql-driver/mysql"
	"os"
	"schoolserver/config"
	"schoolserver/dao/redisDao"
	"schoolserver/logger"
	"xorm.io/xorm"
	"xorm.io/xorm/log"

	"time"
)

var (
	engine        *xorm.Engine
	engineSchool  *xorm.Engine
	engineMessage *xorm.Engine
	logLevel      string
	logPath       string
)

var levelMap = map[string]log.LogLevel{
	"debug": log.LOG_DEBUG,
	"info":  log.LOG_INFO,
	"error": log.LOG_INFO,
}

func NowTimeStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func NowTimeNumStr() string {
	return time.Now().Format("20060102150405")
}

func NowDateStr() string {
	return time.Now().Format("20060102")
}

func InitMysqlDB(config *conf.AllConfig) (err error) {
	if config == nil {
		err = errors.New("initEngine err:config==nil")
		return
	}
	if len(config.Sys.DbLogLevel) > 0 {
		// need more
		logLevel = config.Sys.DbLogLevel
	}
	if len(config.Sys.DbLogFilePath) > 0 {
		logPath = config.Sys.DbLogFilePath
	}

	if &config.MysqlSchool != nil {
		engineSchool, err = initEngine(&config.MysqlSchool)
		if err != nil {
			return
		}
	}

	if &config.Mysql != nil {
		engine, err = initEngine(&config.Mysql)
		if err != nil {
			return
		}
	}

	if &config.MysqlMsg != nil {
		engineMessage, err = initEngine(&config.MysqlMsg)
		if err != nil {
			return
		}
	}

	return
}

func initEngine(configBase *conf.MysqlConfigBase) (engine *xorm.Engine, err error) {
	if configBase == nil {
		err = errors.New("initEngine err: configBase==nil")
		logger.Error(err)
		return
	}

	cfg_ := mysql_.NewConfig()
	cfg_.Net = configBase.Protocol
	addr := configBase.UnixDomain
	if configBase.Protocol == "tcp" {
		addr = fmt.Sprintf("%s:%d", configBase.Host, configBase.Port)
	}
	cfg_.Addr = addr
	cfg_.User = configBase.Username
	cfg_.Passwd = configBase.Password
	cfg_.DBName = configBase.Dbname

	cfg_.AllowOldPasswords = configBase.Params.AllowOldPasswords

	dsn := cfg_.FormatDSN()

	engine, err = xorm.NewEngine("mysql", dsn)
	if err != nil {
		logger.Errorf("Could not initialize DB connection:%v", err)
		panic(err.Error())
		//return errors.New("Open Database Error,Pls See Log!")
	}
	engine.SetMaxOpenConns(configBase.MaxOpenConns)
	engine.SetMaxIdleConns(configBase.MaxIdleConns)
	engine.SetConnMaxLifetime(configBase.ConnMaxLifetime * time.Second)
	engine.ShowSQL(true)

	// set log path and level
	f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.Error(err, logPath)
		return
	}
	engine.SetLogger(log.NewSimpleLogger(f))

	level, in := levelMap[logLevel]
	if in == false {
		return nil, errors.New("wrong db log level")
	}
	engine.Logger().SetLevel(level)
	err = engine.Ping()
	if err != nil {
		logger.Errorf("Could not open  DB connection:%v", err)
		panic(err.Error())
	}
	return
}


func engineWrap(session *xorm.Session, fn func() error) error {
	err := session.Begin()
	defer session.Close()
	if err != nil {
		return err
	}
	if err := fn(); err != nil {
		rollbackErr := session.Rollback()
		if rollbackErr != nil {
			return rollbackErr
		}
		return err
	}
	return session.Commit()
}

func InitRedis(config *conf.RedisOptions) error {
	return redisDao.InitWithConfig(config)
}

func Close() {
	err := engine.Close()
	if err != nil {
		logger.Error("Could not close DB connection:%v", err)
	} else {
		logger.Info("DB closed")
	}

	err = engineMessage.Close()
	if err != nil {
		logger.Error("Could not close DB connection:%v", err)
	} else {
		logger.Info("DB closed")
	}


}
