package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogZap is implementation of HTTPLogger with zap
type LogZap struct {
	logZap *zap.SugaredLogger //zap.NewNop()

}

func NewZapLogger(level string, path string) (*LogZap, error) {
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
		ErrorOutputPaths:  []string{"stderr"},
		DisableStacktrace: true,
	}

	logger, err := conf.Build()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()

	return &LogZap{logger.Sugar()}, nil
}

// RequestLog makes request log
func (logger *LogZap) RequestLog(method string, path string) {
	logger.logZap.Infow("incoming request",
		"method", method,
		"path", path,
	)
}

// Info logs message at info level
func (logger *LogZap) Info(mes string) {
	logger.logZap.Info(mes)
}

func (logger *LogZap) Infof(str string, arg ...any) {
	logger.logZap.Infof(str, arg...)
}
func (logger *LogZap) Errorf(str string, arg ...any) {
	logger.logZap.Errorf(str, arg...)

}

// Error logs message at error level
func (logger *LogZap) Error(mes string) {
	logger.logZap.Error(mes)
}

// Debug logs message at debug level
func (logger *LogZap) Debug(mes string) {
	logger.logZap.Debug(mes)
}

// ResponseLog makes response log
func (logger *LogZap) ResponseLog(status int, size int, duration time.Duration) {
	logger.logZap.Infow("Send response with",
		"status", status,
		"size", size,
		"time", duration.String(),
	)
}
