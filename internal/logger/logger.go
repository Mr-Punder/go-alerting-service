package logger

import (
	"time"

	"go.uber.org/zap"
)

// LogHttpZap is implementation of HttpLogger with zap
type LogHttpZap struct {
	logZap *zap.SugaredLogger //zap.NewNop()

}

func NewLogZap(level string, path string, errorPath string) (*LogHttpZap, error) {
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

	return &LogHttpZap{logger.Sugar()}, nil
}

// RequestLog makes request log
func (logger *LogHttpZap) RequestLog(method string, path string, duration time.Duration) {
	logger.logZap.Infow("incoming request",
		"method", method,
		"path", path,
		"time", duration.String(),
	)
}

func (logger *LogHttpZap) ResponseLog(status int, size int) {
	logger.logZap.Infow("Send response with",
		"status", status,
		"size", size,
	)
}
