// Package that handle chat service in the app
// Currently it will handle chat on the upstore app or through the kakaotalk by hit current API in upstore server
package chat

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"sync"
)

// MessageBatch
type Message struct {
	// Represent ID (not the OID) of sender
	From string `json:"from"`

	// Represent ID (not the OID) of recipient
	To string `json:"to"`

	// Will contains message text
	Msg string `json:"message"`

	// An optional data that needs to be send into the chat set as nil to ignore
	Data interface{} `json:"data,omitempty"`

	// Will contains custom type for mobile
	Type string `json:"type,omitempty"`
}

// KakaoBatch
type Kakao struct {
	// Phone number
	PhoneNumber string `json:"phone_number"`

	// Messages
	Msg string `json:"message"`

	// Template code
	TplCode string `json:"tpl_code"`

	// Button
	Button *string `json:"button"`

	// Item
	Item string `json:"item"`

	// Sending
	Sending string `json:"sending"`

	// Quantity
	Quantity uint `json:"quantity"`
}

// Configuration for chat
type ChatConfig struct {
	ApiToken   string `env-required:"true" env:"UPSTORE_API_BACKSTAGE_TOKEN"`
	BaseApiUrl string `env-required:"true" env:"UPSTORE_API_BASE_URL"`
}

// Chat interface
type Chat interface {
	// Message will handle sending messages through platform
	// Accept string, string, string and interface{} as parameter will return error
	//
	// `from` represent ID (not the OID) of sender
	//
	// `to` represent ID (not the OID) of recipient
	//
	// `message` will contains message text
	//
	// `data` is an optional data that needs to be send into the chat set as nil to ignore
	//
	// e.g.
	// sender := model.Business{
	//		ID: "hangil000",
	// }
	// recipient := model.Business{
	//		ID: "zikeum",
	// }
	// Message(sender.ID, recipient.ID, "Hello zikeum", "1", nil)
	Message(from string, to string, message string, types string, data interface{}) error

	// Message will handle sending messages through platform
	// Accept string, string, string and interface{} as parameter will return error
	//
	// `from` represent ID (not the OID) of sender
	//
	// `to` represent ID (not the OID) of recipient
	//
	// `message` will contains message text
	//
	// `data` is an optional data that needs to be send into the chat set as nil to ignore
	//
	// e.g.
	// sender := model.Business{
	//		ID: "hangil000",
	// }
	// recipient := model.Business{
	//		ID: "zikeum",
	// }
	// sendMessage := []Message{
	//     {
	//		  From: sender.ID,
	//        To: recipient.ID,
	//        Msg: "Hello",
	//        Data: nil
	// 	   }
	// }
	// Message(sender.ID, recipient.ID, "Hello zikeum", nil)
	BatchMessage(messages []Message) error

	// Kakao will handle sending message through kakao talk app
	// Accept string and string will return error
	//
	// e.g.
	// Kakao("0101234567", "Hello", "TP_1234")
	Kakao(phoneNumber string, message string, tplCode string, button *string) error

	// BatchKakao will handle sending message through kakao talk app in batch
	// Accept []Kakao will return error
	//
	// e.g.
	// Kakao([]Kakao{{PhoneNumber: "0101234567890", Msg: "Hello"}})
	BatchKakao(messages []Kakao) error
}

// chat struct will be the base model of other chat struct
type chat struct {
	mtx    sync.Mutex
	client *http.Client
	config *ChatConfig
	body   []byte
	Error  error
}

// Request will do the http request and set error into Error property
// Accept parameter string, string and any will return *chat
//
// Current available method ["GET", "POST"]
func (c *chat) Request(method string, rawUrl string, body any) *chat {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	var reqBody []byte
	var err error

	switch method {
	case http.MethodPost:
		// Marshal request body
		reqBody, err = json.Marshal(&body)
		if err != nil {
			c.Error = err
			return c
		}
		break
	case http.MethodGet:
		break
	default:
		c.Error = errors.New("unavailable method")
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
// Accept parameter any will return error
func (c *chat) Decode(v any) error {
	if c.Error != nil {
		return c.Error
	}

	err := json.Unmarshal(c.body, v)
	if err != nil {
		return err
	}

	return nil
}

// =====
// Private function
// =====

// Do http request
// Accept *http.Request will return []byte, error
func (c *chat) do(req *http.Request) ([]byte, error) {
	// Set default header
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", c.config.ApiToken)

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

// Get the full url based on endpoint url that passed in parameter
// Accept string will return string
//
// e.g.
// GetUrl("/user") will return base_url/user
func (c *chat) getUrl(url string) string {
	return fmt.Sprintf("%s%s", c.config.BaseApiUrl, url)
}
