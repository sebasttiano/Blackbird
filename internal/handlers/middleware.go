package handlers

import (
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"net/http"
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
			status: 0,
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
			"uri", req.RequestURI,
			"method", req.Method,
			"status", responseData.status,
			"duration", duration,
			"size", responseData.size,
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
