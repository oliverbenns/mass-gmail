package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"google.golang.org/api/gmail/v1"
)

type MessageDetails struct {
	SenderName    string `json:"sender_name"`
	SenderAddress string `json:"sender_address"`
	To            string `json:"to"`
	Subject       string `json:"subject"`
	Body          string `json:"body"`
}

func createGmailMessage(m MessageDetails) gmail.Message {
	from := fmt.Sprintf("%s <%s>", m.SenderName, m.SenderAddress)
	parts := []string{
		"From: " + from,
		"To: " + m.To,
		"Subject: " + m.Subject + "\n",
		m.Body,
	}

	str := strings.Join(parts, "\r\n")

	message := gmail.Message{
		Raw: encodeWeb64String([]byte(str)),
	}

	return message
}

func encodeWeb64String(b []byte) string {
	s := base64.URLEncoding.EncodeToString(b)
	return s
}
