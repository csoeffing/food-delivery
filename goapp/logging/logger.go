package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func CreateLogger() *zap.Logger {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/zap.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	level := zap.NewAtomicLevelAt(zap.DebugLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString("\x1b[35mDEBUG\x1b[0m") // Custom Magenta
		case zapcore.InfoLevel:
			enc.AppendString("\x1b[32mINFO\x1b[0m") // Custom Green
		default:
			zapcore.CapitalColorLevelEncoder(l, enc) // Fallback to default color
		}
	}
	//developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	//developmentCfg.EncodeCaller = encoders.ShortColorCallerEncoder(&decorate.Format{Foreground: colors.Yellow})

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	logger := zap.New(core)

	return logger
}

func CreateSugared() *zap.SugaredLogger {
	logger := CreateLogger()

	sugar := logger.Sugar()

	return sugar
}
