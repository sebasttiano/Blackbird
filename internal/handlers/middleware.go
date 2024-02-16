package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/sebasttiano/Blackbird.git/internal/common"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
	"time"
)

// OnlyPostAllowed gets Handler, make validation on Post request and returns also Handler.
func OnlyPostAllowed(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Only POST method allowed
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		next.ServeHTTP(res, req)
	})
}

// WithLogging - middleware that logs request and response params
func WithLogging(next http.Handler) http.Handler {
	logFn := func(res http.ResponseWriter, req *http.Request) {

		start := time.Now()

		responseData := &responseData{
			status: 200,
			size:   0,
		}

		lw := loggingResponseWriter{
			ResponseWriter: res,
			responseData:   responseData,
		}

		next.ServeHTTP(&lw, req)

		duration := time.Since(start)
		sugar := logger.Log.Sugar()
		sugar.Infoln(
			zap.String("uri", req.RequestURI),
			zap.String("method", req.Method),
			zap.Int("status", responseData.status),
			zap.Duration("duration", duration),
			zap.Int("size", responseData.size),
		)

	}
	return http.HandlerFunc(logFn)
}

type (
	// for response data
	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// GzipMiddleware handles compressed with gzip requests and responses
func GzipMiddleware(next http.Handler) http.Handler {

	gzipFn := func(res http.ResponseWriter, req *http.Request) {

		ow := res

		acceptEncoding := req.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := common.NewGZIPWriter(res)
			ow = cw
			defer cw.Close()
		}

		contentEncoding := req.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := common.NewZIPReader(req.Body)
			if err != nil {
				logger.Log.Error("couldn`t decompress request")
				res.WriteHeader(http.StatusInternalServerError)
				return
			}
			req.Body = cr
			defer cr.Close()
		}

		next.ServeHTTP(ow, req)
	}
	return http.HandlerFunc(gzipFn)
}

func CheckSign(key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			hashSHA256 := req.Header.Get("HashSHA256")
			if hashSHA256 == "" {
				next.ServeHTTP(res, req)
				return
			}

			h := hmac.New(sha256.New, []byte(key))

			b, err := io.ReadAll(req.Body)
			if err != nil {
				logger.Log.Error("failed to read request body")
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			if _, err := h.Write(b); err != nil {
				logger.Log.Error("failed to write bytes to hmac")
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			sign := h.Sum(nil)
			headerSign, err := hex.DecodeString(hashSHA256)
			if err != nil {
				logger.Log.Error("failed to decode hashSHA256 header hash")
				res.WriteHeader(http.StatusInternalServerError)
				return
			}

			if !hmac.Equal(sign, headerSign) {
				logger.Log.Error("error: signature validation failed")
				res.WriteHeader(http.StatusBadRequest)
				return
			}

			req.Body = io.NopCloser(bytes.NewReader(b))
			next.ServeHTTP(res, req)
		})
	}
}
