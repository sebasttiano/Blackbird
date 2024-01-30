package handlers

import (
	"github.com/sebasttiano/Blackbird.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpdateMetric(t *testing.T) {

	tests := []struct {
		name         string
		method       string
		expectedCode int
		expectedBody string
		url          string
	}{
		{name: "Check POST /", url: "/", method: http.MethodPost, expectedCode: http.StatusMethodNotAllowed, expectedBody: ""},
		{name: "Check /update/counter/", url: "/update/counter", method: http.MethodPost, expectedCode: http.StatusNotFound, expectedBody: ""},
		{name: "Check /update/gauge", url: "/update/gauge", method: http.MethodPost, expectedCode: http.StatusNotFound, expectedBody: ""},
		{name: "Pass gauge metric #1", url: "/update/gauge/TestMetric/333.3453", method: http.MethodPost, expectedCode: http.StatusOK, expectedBody: ""},
		{name: "Pass gauge metric #2", url: "/update/gauge/TestMetric/133.3453", method: http.MethodPost, expectedCode: http.StatusOK, expectedBody: ""},
		{name: "Pass counter metric #1", url: "/update/counter/TestMetric/10", method: http.MethodPost, expectedCode: http.StatusOK, expectedBody: ""},
		{name: "Pass counter metric #2", url: "/update/counter/TestMetric/20", method: http.MethodPost, expectedCode: http.StatusOK, expectedBody: ""},
		{name: "Get current gauge metric", url: "/value/gauge/TestMetric", method: http.MethodGet, expectedCode: http.StatusOK, expectedBody: "133.3453"},
		{name: "Get current counter metric", url: "/value/counter/TestMetric", method: http.MethodGet, expectedCode: http.StatusOK, expectedBody: "30"},
		{name: "Bad metric type", url: "/update/countere/TestMetric/20", method: http.MethodPost, expectedCode: http.StatusBadRequest, expectedBody: ""},
		{name: "Bad counter value", url: "/update/counter/TestMetric/20ad", method: http.MethodPost, expectedCode: http.StatusBadRequest, expectedBody: ""},
		{name: "Bad gauge value", url: "/update/gauge/TestMetric/aeew", method: http.MethodPost, expectedCode: http.StatusBadRequest, expectedBody: ""},
	}

	views := NewServerViews(storage.NewMemStorage(&storage.StoreSettings{SyncSave: false}))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			router := views.InitRouter()
			router.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			if tt.expectedBody != "" {
				assert.Contains(t, strings.TrimSpace(w.Body.String()), tt.expectedBody, "Содержимое тело ответа не совпадает с ожидаемым")
			}
		})
	}
}
