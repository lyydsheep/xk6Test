package log

import (
	"email/common/enum"
	"email/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

// 使用 zap 作为日志库

var zapLogger *zap.Logger

func InitLogger() {
	// 创建一个适用于生产环境的编码器配置
	encoderConfig := zap.NewProductionEncoderConfig()

	// 设置时间编码方式为ISO8601格式，以提高日志的可读性和国际化
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	// 基于上述配置创建一个JSON编码器，用于生成JSON格式的日志
	fileWriteSyncer, err := getFileLogWriter()
	if err != nil {
		panic(err)
	}
	var cores []zapcore.Core

	switch config.App.Env {
	// 只输出到文件
	case enum.ModeTEST, enum.ModePROD:
		cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(fileWriteSyncer), zap.InfoLevel))

	// 同时输出到控制台和文件
	case enum.ModeDEV:
		cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(fileWriteSyncer), zap.InfoLevel))
		cores = append(cores, zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(os.Stdout), zap.DebugLevel))
	}
	core := zapcore.NewTee(cores...)

	zapLogger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func getFileLogWriter() (zapcore.WriteSyncer, error) {
	// 初始化 lumberjack.Logger 实例，配置日志文件的路径、最大大小、最大年龄，
	lumberJackLogger := &lumberjack.Logger{
		Filename:  config.App.Log.Path,
		MaxSize:   config.App.Log.MaxSize,
		MaxAge:    config.App.Log.MaxAge,
		Compress:  true,
		LocalTime: true,
	}
	// 使用 zapcore.AddSync 将 lumberjack.Logger 转换为 zapcore.WriteSyncer 接口，
	// 以便在 zap 日志库中使用，并返回。
	return zapcore.AddSync(lumberJackLogger), nil
}
