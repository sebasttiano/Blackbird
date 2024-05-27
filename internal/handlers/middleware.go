package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sebasttiano/Blackbird.git/internal/common"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"go.uber.org/zap"
)

// OnlyPostAllowed проверяет метод на POST
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

// WithLogging - логгирует все request запросы и response ответы
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
	// responseData хранит статус коди и размер ответа
	responseData struct {
		status int
		size   int
	}
	// loggingResponseWriter реализует http.ResponseWriter и доп. данные для ответ
	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

// Write проксирует ответ и логгирует его
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

// WriteHeader проксирует заголовки и логгирует код ответа
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

// GzipMiddleware сжимает и распаковывает gzip данные из запроса и ответа
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

// WithRSADecryption decrypts incoming requests body
func WithRSADecryption(priv *rsa.PrivateKey) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		encFn := func(res http.ResponseWriter, req *http.Request) {
			if priv == nil || req.Method == http.MethodGet {
				next.ServeHTTP(res, req)
				return
			}

			b, err := io.ReadAll(req.Body)
			if err != nil {
				logger.Log.Error("failed to read request body")
				http.Error(res, "failed to read request body, check your request", http.StatusBadRequest)
				return
			}

			decrypted, err := common.DecryptRSA(string(b), priv)
			if err != nil {
				logger.Log.Error("failed to decrypt request")
				http.Error(res, "failed to decrypt request, check your request", http.StatusBadRequest)
			}
			req.Body = io.NopCloser(bytes.NewBufferString(decrypted))
			next.ServeHTTP(res, req)
		}
		return http.HandlerFunc(encFn)
	}
}

// CheckSign проверяет цифровую подпись, если есть соответсвующий заголовок
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
				http.Error(res, "failed to read request body, check your request", http.StatusBadRequest)
				return
			}

			if _, errWr := h.Write(b); errWr != nil {
				logger.Log.Error("failed to write bytes to hmac")
				http.Error(res, "internal server error", http.StatusInternalServerError)
				return
			}

			sign := h.Sum(nil)
			headerSign, err := hex.DecodeString(hashSHA256)
			if err != nil {
				logger.Log.Error("failed to decode hashSHA256 header hash")
				http.Error(res, "failed to decode hashSHA256", http.StatusInternalServerError)
				return
			}

			if !hmac.Equal(sign, headerSign) {
				logger.Log.Error("error: signature validation failed")
				http.Error(res, "error: signature validation failed", http.StatusInternalServerError)
				return
			}

			req.Body = io.NopCloser(bytes.NewReader(b))
			next.ServeHTTP(res, req)
		})
	}
}
