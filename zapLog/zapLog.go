package zapLog

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger

func init() {
	logger, _ = func() (*zap.Logger, error) {
		//[yyyy-MM-dd hh:mm:ss,209 5435 ][$level ][$location]|$log_message

		pwd, _ := os.Getwd()
		log, err := zap.Config{
			Encoding: "console", //"json",
			Level:    zap.NewAtomicLevelAt(zapcore.DebugLevel),
			OutputPaths: []string{
				"stdout",
				pwd + "/zapLog.log",
			},
			ErrorOutputPaths: []string{
				"stderr",
			},
			EncoderConfig: zapcore.EncoderConfig{
				MessageKey: "message",

				LevelKey: "level",
				EncodeLevel: func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
					enc.AppendString("[" + level.CapitalString() + " ]")
				},

				TimeKey: "time",
				EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
					enc.AppendString("[" + t.Format("2006-01-02 15:04:05.000") + " ]")
				},

				CallerKey: "caller",
				EncodeCaller: func(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
					enc.AppendString("[" + caller.TrimmedPath() + "]")
				},
			},
		}.Build(zap.AddCallerSkip(1))
		return log, err
	}()
	defer func() {
		logger.Sync()
	}()
}

func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}
