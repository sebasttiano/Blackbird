package common

import (
	"errors"
	"fmt"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type HTTPClientErrors struct {
	ErrConnect error
}

func NewHTTPClientErrors() HTTPClientErrors {
	return HTTPClientErrors{
		ErrConnect: errors.New("client couldn`t connect to server"),
	}
}

// HTTPClient simple client
type HTTPClient struct {
	url          string
	client       *http.Client
	retries      int
	retriesIn    []uint
	ClientErrors HTTPClientErrors
}

func NewHTTPClient(url string, retries int, backoffFactor uint) HTTPClient {
	var ri []uint
	for i := 1; i <= retries; i++ {
		ri = append(ri, backoffFactor*uint(i)-1)
	}
	return HTTPClient{url: url, client: &http.Client{}, retriesIn: ri, retries: retries, ClientErrors: NewHTTPClientErrors()}
}

// Post implements http post requests
func (c HTTPClient) Post(urlSuffix string, body io.Reader, headers map[string]string) (*http.Response, error) {

	r, err := http.NewRequest("POST", c.url+urlSuffix, body)
	if err != nil {
		logger.Log.Debug("failed to make http request", zap.Error(err))
		return nil, err
	}
	for key, value := range headers {
		r.Header.Add(key, value)
	}

	var res *http.Response
	for _, delay := range c.retriesIn {
		res, err = c.client.Do(r)
		if err != nil {
			c.retries -= 1
			logger.Log.Error(fmt.Sprintf("Request to server failed. retrying in %d seconds... Retries left %d\n", delay, c.retries))
			time.Sleep(time.Duration(delay) * time.Second)
			if c.retries == 0 {
				return nil, c.ClientErrors.ErrConnect
			}
		} else {
			break
		}
	}

	return res, nil
}
