package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func NewFile(level, filepath string) (*zap.Logger, error) {
	lvl, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath,
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     7, // days
	})
	enc := zap.NewProductionEncoderConfig()
	enc.TimeKey = "ts"
	enc.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(enc),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), w),
		lvl,
	)
	return zap.New(core, zap.AddCaller()), nil
}