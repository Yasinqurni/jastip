package httphelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jastip-app/config"
	"net/http"
	"sync"
)

type HttpHelper struct {
	mtx         sync.Mutex
	client      *http.Client
	config      *config.Config
	body        []byte
	Error       error
	setHeaderFn func(*http.Request) error
}

func NewHttpHelper(client *http.Client, config *config.Config) HttpHelper {
	return HttpHelper{
		client: client,
		config: config,
	}
}

// Request will do the http request and set error into Error property
// Accept parameter string, string and interface{} will return *HttpHelper
//
// Current available method ["GET", "POST"]
func (c *HttpHelper) Request(method string, rawUrl string, body interface{}) *HttpHelper {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	var reqBody []byte
	var err error

	// Marshal request body
	if body == nil {
		body = map[string]any{}
	}
	reqBody, err = json.Marshal(&body)
	if err != nil {
		c.Error = err
		return c
	}

	// Create new http request
	req, err := http.NewRequest(method, rawUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		c.Error = err
		return c
	}

	// Do the http request
	resBody, err := c.do(req)
	if err != nil {
		c.Error = err
		return c
	}

	c.body = resBody
	return c
}

// Decode will unmarshal response body
// Accept parameter interface{} will return error
func (c *HttpHelper) Decode(v interface{}) error {
	if c.Error != nil {
		return c.Error
	}

	err := json.Unmarshal(c.body, v)
	if err != nil {
		return err
	}

	return nil
}

func (c *HttpHelper) SetHeaderFn(fn func(req *http.Request) error) {
	c.setHeaderFn = fn
}

// =====
// Private function
// =====

// Do http request
// Accept *http.Request will return []byte, error
func (c *HttpHelper) do(req *http.Request) ([]byte, error) {
	// Set default header
	req.Header.Set("Content-Type", "application/json")
	c.setHeaderFn(req)

	// Do the request
	res, err := c.client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	// Read all response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	// Check the status code
	if res.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("failed to do request with %d status code and %v", res.StatusCode, string(body))
	}

	return body, nil
}
