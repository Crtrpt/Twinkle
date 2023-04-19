package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var Logger *zap.Logger
var Sugar *zap.SugaredLogger

func init() {
	logDir := "./log/"
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	errorWriteSyncer := zapcore.AddSync(getLogWriter(logDir+"error.log", 1, 3, 28))
	infoWriteSyncer := zapcore.AddSync(getLogWriter(logDir+"info.log", 1, 3, 28))

	//error写文件
	//error info 写日志
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.ErrorLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		zapcore.NewCore(encoder, errorWriteSyncer, zapcore.ErrorLevel),
		zapcore.NewCore(encoder, infoWriteSyncer, zapcore.InfoLevel),
	)

	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	Sugar = Logger.Sugar()
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	//file, _ := os.OpenFile("./log/error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 644)
	lumberjack := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackup,
		LocalTime:  false,
		Compress:   false,
	}
	return zapcore.AddSync(lumberjack)
}

func Debug(args ...interface{}) {
	Sugar.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	Sugar.Debugf(template, args...)
}

func Info(args ...interface{}) {
	Sugar.Info(args...)
}

func Infof(template string, args ...interface{}) {
	Sugar.Infof(template, args...)
}

func Warn(args ...interface{}) {
	Sugar.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	Sugar.Warnf(template, args...)
}

func Error(args ...interface{}) {
	Sugar.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	Sugar.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	Sugar.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	Sugar.DPanicf(template, args...)
}
