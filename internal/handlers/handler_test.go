package handlers

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name    string
		storage MemStorage
		want    MemStorage
	}{
		{
			name:    "Create New MemStorage",
			storage: NewMemStorage(),
			want: MemStorage{
				gauge:   make(map[string]float64),
				counter: map[string]int64{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.storage)
		})
	}
}

func TestUpdateMetric(t *testing.T) {

	tests := []struct {
		name         string
		method       string
		expectedCode int
		expectedBody string
		url          string
	}{
		{name: "Check /", url: "/", method: http.MethodPost, expectedCode: http.StatusNotFound, expectedBody: ""},
		{name: "Check /update/", url: "/update/", method: http.MethodPost, expectedCode: http.StatusNotFound, expectedBody: ""},
		{name: "Check /update/counter/", url: "/update/counter", method: http.MethodPost, expectedCode: http.StatusNotFound, expectedBody: ""},
		{name: "Check /update/gauge", url: "/update/gauge", method: http.MethodPost, expectedCode: http.StatusNotFound, expectedBody: ""},
		{name: "Pass gauge metric", url: "/update/gauge/TestMetric/333.3453", method: http.MethodPost, expectedCode: http.StatusOK, expectedBody: "333.35"},
		{name: "Pass counter metric", url: "/update/counter/TestMetric/10", method: http.MethodPost, expectedCode: http.StatusOK, expectedBody: "10"},
		{name: "Pass counter metric again", url: "/update/counter/TestMetric/20", method: http.MethodPost, expectedCode: http.StatusOK, expectedBody: "30"},
		{name: "Bad metric type", url: "/update/countere/TestMetric/20", method: http.MethodPost, expectedCode: http.StatusBadRequest, expectedBody: ""},
		{name: "Bad counter value", url: "/update/counter/TestMetric/20ad", method: http.MethodPost, expectedCode: http.StatusBadRequest, expectedBody: ""},
		{name: "Bad gauge value", url: "/update/gauge/TestMetric/aeew", method: http.MethodPost, expectedCode: http.StatusBadRequest, expectedBody: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, tt.url, nil)
			w := httptest.NewRecorder()

			mux := NewMetricHandler()
			mux.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}
