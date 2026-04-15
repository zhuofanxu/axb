package logger

import (
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
)

func getEncoderConfig(conf *Config, allowColor bool) (config zapcore.EncoderConfig) {
	encodeCaller := zapcore.ShortCallerEncoder
	if conf.LogInConsole {
		encodeCaller = zapcore.FullCallerEncoder
	}

	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  conf.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     customTimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   encodeCaller,
	}

	config.EncodeLevel = levelEncoderFromName(conf.EncodeLevel, allowColor)
	return config
}

func levelEncoderFromName(name string, allowColor bool) zapcore.LevelEncoder {
	switch name {
	case "LowercaseLevelEncoder":
		return zapcore.LowercaseLevelEncoder
	case "LowercaseColorLevelEncoder":
		if allowColor {
			return zapcore.LowercaseColorLevelEncoder
		}
		return zapcore.LowercaseLevelEncoder
	case "CapitalLevelEncoder":
		return zapcore.CapitalLevelEncoder
	case "CapitalColorLevelEncoder":
		if allowColor {
			return zapcore.CapitalColorLevelEncoder
		}
		return zapcore.CapitalLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func getEncoder(conf *Config, allowColor bool) zapcore.Encoder {
	if conf.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig(conf, allowColor))
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig(conf, allowColor))
}

func getWriteSyncer(file string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   file, // 日志文件的位置
		MaxSize:    30,   // 在进行切割之前，日志文件的最大大小（以MB为单位）
		MaxBackups: 60,   // 保留旧文件的最大个数
		MaxAge:     7,    // 保留旧文件的最大天数
		Compress:   true, // 是否压缩/归档旧文件
		LocalTime:  true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func getEncoderCore(fileName string, level zapcore.LevelEnabler, conf *Config) (core zapcore.Core) {
	writer := getWriteSyncer(fileName)
	return zapcore.NewCore(getEncoder(conf, false), writer, level)
}

func getConsoleCore(level zapcore.LevelEnabler, conf *Config) (core zapcore.Core) {
	writer := zapcore.AddSync(os.Stdout)
	return zapcore.NewCore(getEncoder(conf, true), writer, level)
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
