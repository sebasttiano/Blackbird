package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sebasttiano/Blackbird.git/internal/repository"
	"github.com/sebasttiano/Blackbird.git/internal/service"
	"github.com/stretchr/testify/assert"
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
		&service.Settings{SyncSave: false, Retries: 1, BackoffFactor: 1},
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

func TestUpdateMetricJSON(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Check POST /update counter",
			method:       http.MethodPost,
			body:         `{"id": "PollCount", "type": "counter", "delta": 33, "value": "124,5"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"delta":33, "id":"PollCount", "type":"counter", "value":0}`,
		},
		{
			name:         "Check POST /update counter2",
			method:       http.MethodPost,
			body:         `{"id": "PollCount", "type": "counter", "delta": 67, "value": "0"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"delta":67, "id":"PollCount", "type":"counter", "value":0}`,
		},
		{
			name:         "Check POST /update gauge",
			method:       http.MethodPost,
			body:         `{"id": "allocMem", "type": "gauge", "delta": 0, "value": "124,5"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id": "allocMem", "type": "gauge", "delta": 0, "value": 0}`,
		},
	}

	views := NewServerViews(service.NewService(
		&service.Settings{SyncSave: false, Retries: 1, BackoffFactor: 1},
		repository.NewMemStorage()))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var metrics models.Metrics
			_ = json.Unmarshal([]byte(tt.body), &metrics)

			jsonValue, err := json.Marshal(metrics)
			assert.NoErrorf(t, err, "Ошибка при сериализации в JSON")

			r := httptest.NewRequest(tt.method, "/update/", bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()

			r.Header.Set("Content-Type", "application/json")

			router := views.InitRouter()
			router.ServeHTTP(w, r)
			assert.Equal(t, tt.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String())
			}
		})
	}
}

func TestUpdateMetricsJSON(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		body         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "OK Check POST /updates",
			method:       http.MethodPost,
			body:         `[{"id": "PollCount", "type": "counter", "delta": 33, "value": "124,5"}, {"id": "allocMem", "type": "gauge", "delta": 0, "value": "124,5"}]`,
			expectedCode: http.StatusOK,
			expectedBody: `[{"id":"PollCount","type":"counter","delta":33,"value":0},{"id":"allocMem","type":"gauge","delta":0,"value":0}]`,
		},
		{
			name:         "NOT OK Check POST /updates",
			method:       http.MethodPost,
			body:         `[{"metrica": "PollCount", "code": "counter", "delta": 33, "value": "124,5"}, {"metrica": "allocMem", "type": "gauge", "delta": 0, "value": "124,5"}]`,
			expectedCode: http.StatusBadRequest,
			expectedBody: `[{"id":"PollCount","type":"counter","delta":33,"value":0},{"id":"allocMem","type":"gauge","delta":0,"value":0}]`,
		},
		{
			name:         "OK Check POST /updates",
			method:       http.MethodPut,
			body:         `[{"id": "PollCount", "type": "counter", "delta": 33, "value": "124,5"}, {"id": "allocMem", "type": "gauge", "delta": 0, "value": "124,5"}]`,
			expectedCode: http.StatusMethodNotAllowed,
			expectedBody: `[{"id":"PollCount","type":"counter","delta":33,"value":0},{"id":"allocMem","type":"gauge","delta":0,"value":0}]`,
		},
	}

	views := NewServerViews(service.NewService(
		&service.Settings{SyncSave: false, Retries: 1, BackoffFactor: 1},
		repository.NewMemStorage()))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var metrics []*models.Metrics
			_ = json.Unmarshal([]byte(tt.body), &metrics)
			jsonValue, err := json.Marshal(metrics)
			assert.NoErrorf(t, err, "Ошибка при сериализации в JSON")

			r := httptest.NewRequest(tt.method, "/updates/", bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()

			r.Header.Set("Content-Type", "application/json")

			router := views.InitRouter()
			router.ServeHTTP(w, r)
			assert.Equal(t, tt.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			if w.Code == http.StatusOK {
				if tt.expectedBody != "" {
					assert.JSONEq(t, tt.expectedBody, w.Body.String())
				}
			}
		})
	}
}

func TestGetMetricJSON(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		body         string
		add          string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "OK counter",
			method:       http.MethodPost,
			body:         `{"id": "PollCount", "type": "counter"}`,
			add:          `{"id": "PollCount", "type": "counter", "delta": 50, "value": "0"}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"PollCount", "type":"counter", "delta":50}`,
		},
		{
			name:         "OK gauge",
			method:       http.MethodPost,
			body:         `{"id": "alloc", "type": "gauge"}`,
			add:          `{"id": "alloc", "type": "gauge", "value": 33.1}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":"alloc", "type":"gauge", "value":33.1}`,
		},
		{
			name:         "NOT OK. NOT FOUND",
			method:       http.MethodPost,
			body:         `{"id": "PollCountNOTFOUND", "type": "counter"}`,
			add:          `{"id": "PollCount", "type": "counter", "delta": 50, "value": "0"}`,
			expectedCode: http.StatusNotFound,
			expectedBody: `{"id":"PollCount", "type":"counter", "delta":50}`,
		},
	}

	views := NewServerViews(service.NewService(
		&service.Settings{SyncSave: false, Retries: 1, BackoffFactor: 1},
		repository.NewMemStorage()))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var metrics models.Metrics
			_ = json.Unmarshal([]byte(tt.add), &metrics)
			jsonValue, err := json.Marshal(metrics)
			assert.NoErrorf(t, err, "Ошибка при сериализации в JSON")

			// Add value
			ra := httptest.NewRequest(http.MethodPost, "/update/", bytes.NewBuffer(jsonValue))
			wa := httptest.NewRecorder()

			// Get value
			r := httptest.NewRequest(tt.method, "/value/", bytes.NewBuffer([]byte(tt.body)))
			w := httptest.NewRecorder()

			r.Header.Set("Content-Type", "application/json")
			ra.Header.Set("Content-Type", "application/json")

			router := views.InitRouter()
			router.ServeHTTP(wa, ra)
			router.ServeHTTP(w, r)
			assert.Equal(t, tt.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			if w.Code == http.StatusOK {
				if tt.expectedBody != "" {
					assert.JSONEq(t, tt.expectedBody, w.Body.String())
				}
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

func Example_sign() {
	res1 := sign("example", "SECRET")
	fmt.Println(res1)

	res2 := sign("lllllllll", "SUPERSECRET")
	fmt.Println(res2)

	// Output:
	// 800b896fe5bb8bce7d8a3d3dae28fbd4e8968cde2271449704645796902aed04
	// 4e54afaf64269c1c0c1a0857c0066393eb14b1c63fc87dd66624dbd8acef6eb8
}
