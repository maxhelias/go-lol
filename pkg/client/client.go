package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	lcu *Client
)

type Client struct {
	Port       int
	AuthToken  string
	BaseURL    string
	httpClient *http.Client
}

func NewClient(port int, token string, LcuCertPath string) *Client {
	lcu = &Client{
		Port:      port,
		AuthToken: token,
		BaseURL:   fmt.Sprintf("https://127.0.0.1:%d", port),
	}

	// Load PEM
	certPEM, err := os.ReadFile(LcuCertPath)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(certPEM)

	lcu.httpClient = &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,
			TLSClientConfig: &tls.Config{
				RootCAs:    certPool,
				MinVersion: tls.VersionTLS12,
			},
		},
	}

	return lcu
}

func Get(url string) ([]byte, error) {
	return req(http.MethodGet, url, nil)
}
func Post(url string, body interface{}) ([]byte, error) {
	return req(http.MethodPost, url, body)
}
func Patch(url string, body interface{}) ([]byte, error) {
	return req(http.MethodPatch, url, body)
}
func Del(url string) ([]byte, error) {
	return req(http.MethodDelete, url, nil)
}

func req(method string, url string, data interface{}) ([]byte, error) {
	var body io.Reader
	if data != nil {
		bts, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		body = bytes.NewReader(bts)
	}
	req, _ := http.NewRequest(method, lcu.BaseURL+url, body)
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte("riot:"+lcu.AuthToken)))
	if req.Body != nil {
		req.Header.Add("ContentType", "application/json")
	}

	resp, err := lcu.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)

		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
