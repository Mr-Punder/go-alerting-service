package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogHTTPZap is implementation of HTTPLogger with zap
type LogHTTPZap struct {
	logZap *zap.SugaredLogger //zap.NewNop()

}

func NewLogZap(level string, path string, errorPath string) (*LogHTTPZap, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	conf := zap.Config{
		Level:             lvl,
		Development:       true,
		Encoding:          "console",
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{path},
		ErrorOutputPaths:  []string{errorPath},
		DisableStacktrace: true,
	}

	logger, err := conf.Build()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()

	return &LogHTTPZap{logger.Sugar()}, nil
}

// RequestLog makes request log
func (logger *LogHTTPZap) RequestLog(method string, path string) {
	logger.logZap.Infow("incoming request",
		"method", method,
		"path", path,
	)
}

// Info logs message at info level
func (logger *LogHTTPZap) Info(mes string) {
	logger.logZap.Info(mes)
}

// Error logs message at error level
func (logger *LogHTTPZap) Error(mes string) {
	logger.logZap.Error(mes)
}

// Debug logs message at debug level
func (logger *LogHTTPZap) Debug(mes string) {
	logger.logZap.Debug(mes)
}

// ResponseLog makes response log
func (logger *LogHTTPZap) ResponseLog(status int, size int, duration time.Duration) {
	logger.logZap.Infow("Send response with",
		"status", status,
		"size", size,
		"time", duration.String(),
	)
}
