package logger

import (
	"time"

	"go.uber.org/zap"
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

	conf := zap.Config{
		Level:            lvl,
		Development:      true,
		Encoding:         "json",
		EncoderConfig:    zap.NewDevelopmentEncoderConfig(),
		OutputPaths:      []string{path},
		ErrorOutputPaths: []string{errorPath},
	}

	logger, err := conf.Build()
	if err != nil {
		return nil, err
	}
	defer logger.Sync()

	return &LogHTTPZap{logger.Sugar()}, nil
}

// RequestLog makes request log
func (logger *LogHTTPZap) RequestLog(method string, path string, duration time.Duration) {
	logger.logZap.Infow("incoming request",
		"method", method,
		"path", path,
		"time", duration.String(),
	)
}

func (logger *LogHTTPZap) ResponseLog(status int, size int) {
	logger.logZap.Infow("Send response with",
		"status", status,
		"size", size,
	)
}
