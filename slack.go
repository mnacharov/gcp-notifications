package main

import (
	"encoding/json"
	"bytes"
	"net/http"
)


type SlackMessage struct {
	Text        string `json:"text"`
	Markdown    bool   `json:"mrkdwn"`
}

func sendSlack(url string, msg string) error {
	var m SlackMessage
	m.Text = msg
	m.Markdown = true
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	_, err = http.DefaultClient.Do(req)
	return err
}
