package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"

	"github.com/Mr-Punder/go-alerting-service/internal/logger"
)

type hashWriter struct {
	http.ResponseWriter
	hw bytes.Buffer
}

func NewHashWriter(w http.ResponseWriter) *hashWriter {
	return &hashWriter{
		ResponseWriter: w,
		hw:             bytes.Buffer{},
	}
}

func (hw *hashWriter) Write(b []byte) (int, error) {
	hw.ResponseWriter.WriteHeader(200)
	return hw.hw.Write(b)
}
func (hw *hashWriter) Header() http.Header {
	return hw.ResponseWriter.Header()
}

func (hw *hashWriter) WriteHeader(StatusCode int) {

	hw.ResponseWriter.WriteHeader(StatusCode)
}

type HashSum struct {
	log logger.Logger
	key string
}

func NewHashSum(key string, log logger.Logger) *HashSum {
	return &HashSum{
		key: key,
		log: log,
	}
}

func (hs *HashSum) HashSummHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqHash := r.Header.Get("HashSHA256")
		if reqHash != "" {
			hs.log.Info("Hash sha256 detected")
			body, err := io.ReadAll(r.Body)
			if err != nil {
				hs.log.Errorf("Can not read request body %s", err)
				http.Error(w, "Can not read request body", http.StatusInternalServerError)
			}
			h := hmac.New(sha256.New, []byte(hs.key))
			h.Write(body)
			hash := hex.EncodeToString(h.Sum(nil))
			hs.log.Infof("Calculated hash: %s", hash)
			hs.log.Infof("Recieved hash: %s", reqHash)

			if hash != reqHash {
				hs.log.Error("Hash doesn't match")
				http.Error(w, "Hash doesn't match", http.StatusBadRequest)
				return
			}
			hs.log.Info("Hash is OK")

			newBuffer := bytes.NewBuffer(body)
			r.Body = io.NopCloser(newBuffer)
		} else {
			hs.log.Error("Hash has not detected")
		}

		if hs.key != "" {
			hashWriter := NewHashWriter(w)
			next.ServeHTTP(hashWriter, r)
			respBody := hashWriter.hw.Bytes()
			h := hmac.New(sha256.New, []byte(hs.key))
			h.Write(respBody)
			hash := hex.EncodeToString(h.Sum(nil))
			w.Header().Set("HashSHA256", hash)
			hs.log.Infof("Resp body %s", string(respBody))
			w.Write(respBody)
		} else {
			next.ServeHTTP(w, r)

		}
		hs.log.Info("Request served from HashHandler")
	})
}
