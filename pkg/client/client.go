package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

var (
	lcu *Client
)

type Option func(o *opt) error

func WithPort(port int) Option {
	return func(o *opt) error {
		o.port = port

		return nil
	}
}

func WithAuthToken(token string) Option {
	return func(o *opt) error {
		o.authToken = token

		return nil
	}
}

func WithLcuCertPath(path string) Option {
	return func(o *opt) error {
		o.lcuCertPath = path

		return nil
	}
}

type opt struct {
	port        int
	authToken   string
	lcuCertPath string
}

type Client struct {
	Port       int
	AuthToken  string
	BaseURL    string
	httpClient *http.Client
}

func NewClient(options ...Option) (*Client, error) {
	opt := &opt{}

	for _, o := range options {
		if err := o(opt); err != nil {
			return nil, err
		}
	}

	lcu = &Client{
		Port:      opt.port,
		AuthToken: opt.authToken,
		BaseURL:   fmt.Sprintf("https://127.0.0.1:%d", opt.port),
	}

	var config *tls.Config
	if opt.lcuCertPath != "" {
		certPEM, err := os.ReadFile(opt.lcuCertPath)
		if err != nil {
			panic(err)
		}

		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(certPEM)
		config = &tls.Config{
			RootCAs:    certPool,
			MinVersion: tls.VersionTLS12,
		}
	} else {
		config = &tls.Config{
			InsecureSkipVerify: true,
		}
	}

	lcu.httpClient = &http.Client{
		Transport: &http.Transport{
			ForceAttemptHTTP2: true,
			TLSClientConfig:   config,
		},
	}

	return lcu, nil
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
	req.SetBasicAuth("riot", lcu.AuthToken)
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
