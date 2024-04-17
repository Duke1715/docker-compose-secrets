package client

import (
	"github.com/go-resty/resty/v2"
)

const HeaderVaultTokenName = "X-Vault-Token"

type HttpClient struct {
	client *resty.Client
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		client: resty.New(),
	}
}

func (h *HttpClient) Get(url string, headers map[string]string) ([]byte, error) {
	resp, err := h.client.R().
		SetHeaders(headers).
		Get(url)
	if err != nil {
		return nil, err
	}

	return resp.Body(), nil
}

func (h *HttpClient) BuildHeaders(customHeaders map[string]string) map[string]string {
	headers := make(map[string]string)

	headers["Accept"] = "application/json"
	headers["Content-Type"] = "application/json"

	for name, value := range customHeaders {
		headers[name] = value
	}

	return headers
}
