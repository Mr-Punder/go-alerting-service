package gzipcomp

import (
	"compress/gzip"
	"io"
	"net/http"
)

type GzipCompressWiter struct {
	http.ResponseWriter
	zw *gzip.Writer
}

func NewGzipCompressWriter(w http.ResponseWriter) *GzipCompressWiter {
	return &GzipCompressWiter{
		ResponseWriter: w,
		zw:             gzip.NewWriter(w),
	}
}

func NewEmptyGzipCompressWriter() *GzipCompressWiter {
	return &GzipCompressWiter{}
}

func (gw *GzipCompressWiter) Name() string {
	return "gzip"
}

func (gw *GzipCompressWiter) SetResponseWriter(w http.ResponseWriter) {
	gw.ResponseWriter = w
}

func (gw *GzipCompressWiter) Header() http.Header {
	return gw.ResponseWriter.Header()
}

func (gw *GzipCompressWiter) Write(b []byte) (int, error) {
	gw.WriteHeader(200)
	return gw.zw.Write(b)
}

func (gw *GzipCompressWiter) WriteHeader(StatusCode int) {
	if StatusCode < 300 {
		gw.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	}
	gw.ResponseWriter.WriteHeader(StatusCode)
}

func (gw *GzipCompressWiter) Close() error {
	return gw.zw.Close()
}

type GzipCompressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func NewGzipCompressReader(r io.ReadCloser) (*GzipCompressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &GzipCompressReader{
		r:  r,
		zr: zr,
	}, nil
}

func NewEmptyGzipCompressReader() *GzipCompressReader {

	return &GzipCompressReader{}
}

func (gr *GzipCompressReader) SetReader(r io.ReadCloser) error {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	gr.r = r
	gr.zr = zr
	return nil
}
func (gw *GzipCompressReader) Name() string {
	return "gzip"
}

func (gr *GzipCompressReader) Read(b []byte) (int, error) {
	return gr.zr.Read(b)
}

func (gr *GzipCompressReader) Close() error {
	if err := gr.r.Close(); err != nil {
		return err
	}
	return gr.zr.Close()
}
