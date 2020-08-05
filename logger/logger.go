package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	conf "schoolserver/config"
)

var (
	allLogger     *zap.SugaredLogger
	gatewayLogger *zap.SugaredLogger
	redisLogger   *zap.SugaredLogger
	httpLogger    *zap.SugaredLogger
)

func init() {
	//PprofLogger = NewLogger("./logs/pprofLogger.log", zapcore.InfoLevel, 128, 12, 7, true, "pprof")
	//MonitorLogger = NewLogger("./logs/monitorLogger.log", zapcore.InfoLevel, 128, 12, 7, true, "monitor")
	//GatewayLogger = NewLogger("./logs/gateway.log", zapcore.DebugLevel, 128, 30, 7, true, "Gateway")
}

var levelMap = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func getLoggerLevel(lvl string) zapcore.Level {
	if level, ok := levelMap[lvl]; ok {
		return level
	}
	return zapcore.InfoLevel
}

func InitLog(config *conf.AllConfig) {
	filePath := config.Sys.LogFileName
	level := getLoggerLevel(config.Sys.LogLevel)
	allLogger = NewLogger(filePath, level, 100, 5, 365, true, "OPENSERVER").Sugar()
}

func Debug(args ...interface{}) {
	allLogger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	allLogger.Debugf(template, args...)
}

func Info(args ...interface{}) {
	allLogger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	allLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	allLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	allLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	allLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	allLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	allLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	allLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	allLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	allLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	allLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	allLogger.Fatalf(template, args...)
}
