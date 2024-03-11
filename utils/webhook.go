package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Field struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type Embed struct {
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Fields      []*Field  `json:"fields,omitempty"`
	Color       int       `json:"color,omitempty"`
	Timestamp   time.Time `json:"timestamp,omitempty"`
}

type Message struct {
	Username  string   `json:"username,omitempty"`
	AvatarURL string   `json:"avatar_url,omitempty"`
	Content   string   `json:"content,omitempty"`
	Embeds    []*Embed `json:"embeds,omitempty"`
}

func SendWebhook(url string, message *Message) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Fatalln("erro to convert message to JSON", err)
		return err
	}

	resp, err := http.Post(url, "applications/json", bytes.NewBuffer(messageBytes))
	if err != nil {
		log.Fatalln("cannot send the discord webhook log", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Fatalln("erro at server response", err)
		return err
	}
	return nil
}
