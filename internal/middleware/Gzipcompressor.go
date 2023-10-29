package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Mr-Punder/go-alerting-service/internal/gzipcomp"
)

type middlewareLogger interface {
	Info(mes string)
	Error(mess string)
	Debug(mess string)
	Infof(str string, arg ...any)
}

type GzipCompressor struct {
	log middlewareLogger
}

func NewGzipCompressor(log middlewareLogger) *GzipCompressor {
	return &GzipCompressor{
		log: log,
	}
}

func (c *GzipCompressor) CompressHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c.log.Info("Entered compressor")
		ow := w

		headers := r.Header
		c.log.Info(fmt.Sprintf("Headers:  %v", headers))

		contentEncoding := r.Header.Get("Content-Encoding")
		c.log.Info(fmt.Sprintf("Content-Encoding = %s", contentEncoding))
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			c.log.Info(fmt.Sprintf("Detected %s compression", "gzip"))

			var err error
			r.Body, err = gzipcomp.NewGzipCompressReader(r.Body)
			if err != nil {
				c.log.Error(fmt.Sprintf("Error setting read buffer for %s compressor", "gzip"))
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			defer r.Body.Close()

		}

		accepEncoding := r.Header.Values("Accept-Encoding")
		c.log.Info(fmt.Sprintf("Accept-Encoding: %v", accepEncoding))

		supportGzip := false

		for _, value := range accepEncoding {
			if strings.Contains(value, "gzip") {
				supportGzip = true
				break
			}
		}

		if supportGzip {
			c.log.Info("Detected gzip support")
			cw := gzipcomp.NewGzipCompressWriter(w)

			ow = cw

			defer cw.Close()

		}
		next.ServeHTTP(ow, r)
		c.log.Info("request served from GzipCompressor")

	})
}
