package logger

import (
	"net/http"
	"os"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"moul.io/http2curl"
)

type LoggerZap struct {
	*zap.Logger
}

func NewLogger() *LoggerZap {
	hook := lumberjack.Logger{
		Filename:   "./logs/app.log",
		MaxSize:    1,
		MaxAge:     2,
		MaxBackups: 10,
		Compress:   true,
	}

	encoder := getEncoderLog()

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(os.Stdout),
			zapcore.AddSync(&hook),
		),
		zapcore.InfoLevel,
	)

	return &LoggerZap{zap.New(core, zap.AddCaller())}
}

func getEncoderLog() zapcore.Encoder {
	encodeConfig := zap.NewProductionEncoderConfig()

	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	encodeConfig.TimeKey = "time"

	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder

	return zapcore.NewJSONEncoder(encodeConfig)
}

func (l *LoggerZap) LogPartnerCall(req *http.Request, responseBody string, duration time.Duration) {
	command, _ := http2curl.GetCurlCommand(req)

	l.Info("PARTNER_API_CALL",
		zap.String("curl", command.String()),
		zap.String("response", responseBody),
		zap.Duration("latency", duration),
	)
}
