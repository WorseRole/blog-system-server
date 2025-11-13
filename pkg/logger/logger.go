package logger

import (
	"blog-system-server/internal/config"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger(cfg *config.Config) error {
	// 设置日志级别
	logLevel := zap.InfoLevel
	switch cfg.LogLevel {
	case "debug":
		logLevel = zap.DebugLevel
	case "info":
		logLevel = zap.InfoLevel
	case "warn":
		logLevel = zap.WarnLevel
	case "error":
		logLevel = zap.ErrorLevel
	}

	// 创建日志目录
	if err := os.MkdirAll("logs", 0755); err != nil {
		return err
	}

	// 配置输出到文件和控制台
	fileWriter := getLogWriter("logs/app.log")
	consoleWriter := zapcore.Lock(os.Stdout)

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 控制台带颜色
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建核心
	core := zapcore.NewTee(
		// 输出到文件（JSON格式）
		zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			fileWriter,
			logLevel,
		),
		// 输出到控制台（带颜色）
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			consoleWriter,
			logLevel,
		),
	)

	Logger = zap.New(core, zap.AddCaller())
	zap.ReplaceGlobals(Logger) // 替换全局logger

	return nil
}

// getLogWriter 创建日志文件writer 也可以在这里调整日志可以用日期进行分割日志
func getLogWriter(filename string) zapcore.WriteSyncer {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// 如果文件打开失败，fallback到标准错误输出
		return zapcore.Lock(os.Stderr)
	}
	return zapcore.Lock(file)
}

// 便捷方法
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}
