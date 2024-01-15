package chat

import (
	"net/http"
)

// Chat on the upstore chat
type UpstoreChat struct {
	*chat
}

func NewUpstoreChat(client *http.Client, cfg *ChatConfig) Chat {
	return &UpstoreChat{
		chat: &chat{
			client: client,
			config: cfg,
		},
	}
}

func (c *UpstoreChat) Message(from string, to string, message string, types string, data interface{}) error {
	body := []Message{
		{
			From: from,
			To:   to,
			Msg:  message,
			Data: data,
			Type: types,
		},
	}
	if err := c.Request(http.MethodPost, c.getUrl("/chat/message"), body).Error; err != nil {
		return err
	}

	return nil
}

func (c *UpstoreChat) BatchMessage(messages []Message) error {
	if err := c.Request(http.MethodPost, c.getUrl("/chat/message"), messages).Error; err != nil {
		return err
	}

	return nil
}

func (c *UpstoreChat) Kakao(phoneNumber string, message string, tplCode string, button *string) error {
	body := []Kakao{
		{
			PhoneNumber: phoneNumber,
			Msg:         message,
			TplCode:     tplCode,
			Button:      button,
		},
	}
	if err := c.Request(http.MethodPost, c.getUrl("/chat/kakao"), body).Error; err != nil {
		return err
	}

	return nil
}

func (c *UpstoreChat) BatchKakao(messages []Kakao) error {
	if err := c.Request(http.MethodPost, c.getUrl("/chat/kakao"), messages).Error; err != nil {
		return err
	}

	return nil
}
