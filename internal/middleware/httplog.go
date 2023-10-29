package middleware

import (
	"fmt"
	"net/http"
	"time"
)

type httpLogger interface {
	Info(mes string)
	Error(mes string)
	Debug(mes string)
	RequestLog(method string, path string)
	ResponseLog(status int, size int, duration time.Duration)
}

type responseData struct {
	status int
	size   int
}

type HttpLogger struct {
	log httpLogger
}

func NewHttpLoger(logger httpLogger) *HttpLogger {
	return &HttpLogger{logger}
}

type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {

	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {

	r.responseData.status = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (l *HttpLogger) HttpLogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := r.Header

		method := r.Method
		path := r.RequestURI
		l.log.RequestLog(method, path)
		l.log.Info(fmt.Sprintf("Headers:  %v", headers))

		start := time.Now()

		resD := &responseData{}

		lw := &loggingResponseWriter{
			ResponseWriter: w,
			responseData:   resD,
		}

		next.ServeHTTP(lw, r)

		l.log.Info("request served from HttpLogger")
		duration := time.Since(start)

		l.log.ResponseLog(resD.status, resD.size, duration)

	})
}
