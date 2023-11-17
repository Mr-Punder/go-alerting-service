package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type LogHTTPLogrus struct {
	log *logrus.Logger
}

func NewLogrusLogger(level string, path string) (*LogHTTPLogrus, error) {
	logger := logrus.New()

	var lvl logrus.Level
	err := lvl.UnmarshalText([]byte(level))
	if err != nil {
		return nil, err
	}
	logger.SetLevel(lvl)

	if path == "stdout" {
		logger.SetOutput(os.Stdout)
	} else {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			return nil, err
		}
		logger.SetOutput(file)

	}

	return &LogHTTPLogrus{log: logger}, nil
}

func (logger *LogHTTPLogrus) Info(mes string) {
	ctime := time.Now().Format("2006-01-02 15:04:05.99999999")
	logger.log.Info(fmt.Sprintf("%s	   %s", ctime, mes))
}

func (logger *LogHTTPLogrus) Error(mes string) {
	ctime := time.Now().Format("2006-01-02 15:04:05.99999999")
	logger.log.Error(fmt.Sprintf("%s	   ERROR %s", ctime, mes))
}

func (logger *LogHTTPLogrus) Debug(mes string) {
	ctime := time.Now().Format("2006-01-02 15:04:05.99999999")
	logger.log.Debug(fmt.Sprintf("%s	   %s", ctime, mes))
}

func (logger *LogHTTPLogrus) Infof(str string, arg ...any) {
	ctime := time.Now().Format("2006-01-02 15:04:05.99999999")
	msg := fmt.Sprintf("%s\t   "+str, append([]interface{}{ctime}, arg...)...)
	logger.log.Info(msg)
}
func (logger *LogHTTPLogrus) Errorf(str string, arg ...any) {
	ctime := time.Now().Format("2006-01-02 15:04:05.99999999")
	msg := fmt.Sprintf("%s	   ERROR "+str, append([]interface{}{ctime}, arg...)...)
	logger.log.Error(msg)
}

func (logger *LogHTTPLogrus) RequestLog(method string, path string) {
	ctime := time.Now().Format("2006-01-02 15:04:05.99999999")
	logger.log.WithFields(logrus.Fields{
		"method": method,
		"path":   path,
	}).Info(fmt.Sprintf("%s	   New incomig request", ctime))
}

func (logger *LogHTTPLogrus) ResponseLog(status int, size int, duration time.Duration) {
	ctime := time.Now().Format("2006-01-02 15:04:05.99999999")
	logger.log.WithFields(logrus.Fields{
		"status":     status,
		"size":       size,
		"spent time": duration.String(),
	}).Info(fmt.Sprintf("%s	   Send response with", ctime))
}
