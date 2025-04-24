package nodes

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/osmoview/flownodes"
)

type TelegramChatBot struct {
	URL string
}

type TelegramChatBotMessage struct {
	Text string `json:"text"`
	Chat int64  `json:"chat"`
}

func (msg TelegramChatBotMessage) IsValid() bool {
	return msg.Text != "" && msg.Chat != 0
}

func (node *TelegramChatBot) Run(msg TelegramChatBotMessage) (flownodes.Data, error) {
	if node.URL == "" {
		return nil, errors.New("URL is empty")
	}

	// if has input data, then send message
	if msg.IsValid() {
		log.Printf("send message: %+v", msg)
		data, err := json.Marshal(msg)
		if err != nil {
			return nil, err
		}

		if _, err := http.Post(node.URL, "application/json", bytes.NewReader(data)); err != nil {
			return nil, err
		}

		return nil, nil
	}

	resp, err := http.Get(node.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// no data no response
	if resp.StatusCode == 204 {
		return nil, nil
	}

	userMsg := map[string]any{}
	if err := json.NewDecoder(resp.Body).Decode(&userMsg); err != nil {
		return nil, err
	}

	log.Printf("get message: %+v", userMsg)

	return userMsg, nil
}
