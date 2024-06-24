package agent

import (
	"bytes"
	"crypto/hmac"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/sebasttiano/Blackbird.git/internal/common"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"github.com/sebasttiano/Blackbird.git/internal/models"
	"go.uber.org/zap"
)

// HTTPSender реализующий интерфейс Sender, отправляет на REST API
type HTTPSender struct {
	client    common.HTTPClient
	signKey   string
	publicKey *rsa.PublicKey
	XRealIP   string
}

// SendToRepo собирает из каналов метрики, формирует и шлет http запрос в репозиторий
func (h *HTTPSender) SendToRepo(jobsMetrics <-chan MetricsSet, jobsGMetrics <-chan GopsutilMetricsSet) error {
	var metric MetricsSet
	var metricG GopsutilMetricsSet
	var metrics models.Metrics
	var metricsBatch []models.Metrics
	var value reflect.Value

	select {
	case metric = <-jobsMetrics:
		value = reflect.ValueOf(metric)
	case metricG = <-jobsGMetrics:
		value = reflect.ValueOf(metricG)
	}
	numFields := value.NumField()
	structType := value.Type()

	for i := 0; i < numFields; i++ {
		field := structType.Field(i)
		fieldValue := value.Field(i)
		metrics.ID = field.Name

		if fieldValue.CanInt() {
			counterVal := fieldValue.Int()
			metrics.Delta = &counterVal
			metrics.MType = "counter"
		} else {
			gaugeVal := fieldValue.Float()
			metrics.Value = &gaugeVal
			metrics.MType = "gauge"
		}

		metricsBatch = append(metricsBatch, metrics)
	}

	if len(metricsBatch) > 0 {
		// Make an HTTP post request
		reqBody, err := json.Marshal(metricsBatch)
		if err != nil {
			logger.Log.Error("couldn`t serialize to json", zap.Error(err))
			return fmt.Errorf("%w: %v", ErrSendToRepo, err)

		}

		compressedData, err := common.Compress(reqBody)
		if err != nil {
			logger.Log.Error("failed to compress data to gzip", zap.Error(err))
			return fmt.Errorf("%w: %v", ErrSendToRepo, err)
		}

		headers := map[string]string{"Content-Type": "application/json", "Content-Encoding": "gzip", "X-Real-IP": h.XRealIP}
		if h.signKey != "" {
			data := *compressedData
			h := hmac.New(sha256.New, []byte(h.signKey))
			if _, errWr := h.Write(data.Bytes()); errWr != nil {
				logger.Log.Error("failed to create hmac signature")
				return fmt.Errorf("%w: %v", ErrSendToRepo, err)
			}
			dst := h.Sum(nil)
			logger.Log.Info("create hmac signature")
			headers["HashSHA256"] = hex.EncodeToString(dst)
		}

		if h.publicKey != nil {
			encrypted, err := common.EncryptRSA(compressedData.String(), h.publicKey)
			if err != nil {
				logger.Log.Error("couldn`t encrypt json data", zap.Error(err))
			}
			compressedData = bytes.NewBuffer([]byte(encrypted))
		}

		res, err := h.client.Post("/updates/", compressedData, headers)
		if err != nil {
			logger.Log.Error(fmt.Sprintf("couldn`t send metrics batch of length %d", len(metricsBatch)), zap.Error(err))
			return fmt.Errorf("%w: %v", ErrSendToRepo, err)
		}
		answer, _ := io.ReadAll(res.Body)
		res.Body.Close()

		if res.StatusCode != 200 {
			logger.Log.Error(fmt.Sprintf("error: server return code %d: message: %s", res.StatusCode, answer))
			return fmt.Errorf("%w: %v", ErrSendToRepo, err)
		}
		logger.Log.Info("send metrics to repository server successfully.")
	}
	return nil
}
