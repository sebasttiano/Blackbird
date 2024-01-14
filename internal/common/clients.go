package common

import (
	"errors"
	"github.com/sebasttiano/Blackbird.git/internal/logger"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strings"
)

// HTTPClient simple client
type HTTPClient struct {
	url string
}

func NewHTTPClient(url string) HTTPClient {
	return HTTPClient{url: url}
}

// Post implements http post requests
func (c HTTPClient) Post(urlSuffix string, body io.Reader, headers []string) (*http.Response, error) {

	r, err := http.NewRequest("POST", c.url+urlSuffix, body)
	if err != nil {
		logger.Log.Debug("failed to make http request", zap.Error(err))
		return nil, err
	}
	for _, header := range headers {
		if header != "" {
			splitHeader := strings.Split(header, ":")
			if len(splitHeader) == 2 {
				r.Header.Add(splitHeader[0], splitHeader[1])
			} else {
				return nil, errors.New("error: check passed header,  it should be in the format '<Name>: <Value>'")
			}
		}
	}
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return nil, err
	}

	return res, nil
}
