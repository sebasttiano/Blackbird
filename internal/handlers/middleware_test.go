package handlers

import (
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sebasttiano/Blackbird.git/internal/repository"
	"github.com/sebasttiano/Blackbird.git/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			})))
			mux.ServeHTTP(w, r)

			assert.Equal(t, tt.expectedCode, w.Code, "Код ответа не совпадает с ожидаемым")
			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, w.Body.String(), "Тело ответа не совпадает с ожидаемым")
			}
		})
	}
}

func TestGzipMiddleware(t *testing.T) {
	views := NewServerViews(service.NewService(&service.ServiceSettings{SyncSave: false}, repository.NewMemStorage()))
	srv := httptest.NewServer(views.InitRouter())
	defer srv.Close()

	tests := []struct {
		name         string
		requestURL   string
		method       string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Check gzipped /update/counter",
			requestURL:   "/update/counter/TestMetric/10",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:         "Check gzipped /update/counter2",
			requestURL:   "/update/counter/TestMetric/20",
			method:       http.MethodPost,
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:         "Check gzipped /value/counter",
			requestURL:   "/value/counter/TestMetric",
			method:       http.MethodGet,
			expectedCode: http.StatusOK,
			expectedBody: "30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(tt.method, srv.URL+tt.requestURL, nil)
			r.RequestURI = ""
			r.Header.Set("Accept-Encoding", "gzip")

			resp, err := http.DefaultClient.Do(r)
			status := resp.StatusCode
			assert.Equal(t, tt.expectedCode, status, "Код ответа не совпадает с ожидаемым")
			require.NoError(t, err)
			defer resp.Body.Close()

			_, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
		})
	}
}

func TestCheckTrustedSubnet(t *testing.T) {

	tests := []struct {
		name          string
		requestURL    string
		trustedSubnet *net.IPNet
		remoteAddr    string
		expectedCode  int
		expectedBody  string
	}{
		{
			name:          "Check ipv4 ok",
			requestURL:    "/update/counter/TestMetric/10",
			trustedSubnet: &net.IPNet{IP: []byte("192.168.1.0"), Mask: net.IPv4Mask(255, 255, 255, 0)},
			remoteAddr:    "192.168.1.100",
			expectedCode:  http.StatusOK,
			expectedBody:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			views := NewServerViews(service.NewService(&service.ServiceSettings{SyncSave: false}, repository.NewMemStorage()))
			views.TrustedSubnet = tt.trustedSubnet
			srv := httptest.NewServer(views.InitRouter())
			defer srv.Close()

			r := httptest.NewRequest(http.MethodPost, srv.URL+tt.requestURL, nil)
			r.RemoteAddr = tt.remoteAddr
			r.RequestURI = ""

			resp, err := http.DefaultClient.Do(r)
			status := resp.StatusCode
			assert.Equal(t, tt.expectedCode, status, "Код ответа не совпадает с ожидаемым")
			require.NoError(t, err)
			defer resp.Body.Close()

			_, err = io.ReadAll(resp.Body)
			require.NoError(t, err)
		})
	}
}
