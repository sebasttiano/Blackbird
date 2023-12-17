package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOnlyPostAllowed(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		expectedBody string
		expectedCode int
	}{
		{method: http.MethodGet, expectedCode: http.StatusMethodNotAllowed, expectedBody: ""},
		{method: http.MethodPut, expectedCode: http.StatusMethodNotAllowed, expectedBody: ""},
		{method: http.MethodDelete, expectedCode: http.StatusMethodNotAllowed, expectedBody: ""},
		{method: http.MethodPost, expectedCode: http.StatusOK, expectedBody: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, "/", nil)
			w := httptest.NewRecorder()

			mux := http.NewServeMux()
			mux.Handle("/", OnlyPostAllowed(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				return
			})))
			mux.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}
