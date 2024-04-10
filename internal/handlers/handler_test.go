package handlers

import (
	"github.com/sebasttiano/Blackbird.git/internal/repository"
	"github.com/sebasttiano/Blackbird.git/internal/service"
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

	views := NewServerViews(service.NewService(
		&service.ServiceSettings{SyncSave: false, Retries: 1, BackoffFactor: 1},
		repository.NewMemStorage()))

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

func BenchmarkHandler_sign(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		sign("182f417b-f260-4b52-ab63-4a74cb7f0555", "oogheemeiS3ailiemaP1eeco9sodai")
	}
}

func Test_sign(t *testing.T) {

	tests := []struct {
		name  string
		value any
		key   string
		want  string
	}{
		{name: "simple sign", value: "test", key: "SECRET", want: "c476367560b2ebfba7b3b1e5b5ef1aa922ee3135db5b4811e94963b6e6ab4ff7"},
		{name: "json", value: "{'id': 'Petrov', 'age': '31'}", key: "PASSPORT", want: "642c92a040538ce901501c66dee5045828bd19448d8dc3a15722ef924ac50566"},
		{name: "uuid", value: "182f417b-f260-4b52-ab63-4a74cb7f0555", key: "UUID", want: "98694e72ab82c15fd51e6e90953ec34a36bacf4f5106005fb6a4ac791470f673"},
		{name: "hard secret key", value: "FOOBAR", key: "oogheemeiS3ailiemaP1eeco9sodai", want: "3a22ea6d51870750131cc1d97ac3f372e93265426e7686cc5b9f5b61b7887abb"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, sign(tt.value, tt.key), "sign(%v, %v)", tt.value, tt.key)
		})
	}
}
