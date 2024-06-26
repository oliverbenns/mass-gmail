package main

import (
	"context"
	"log"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func main() {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client, err := getClient(ctx, config)
	if err != nil {
		log.Fatalf("Unable to get client: %v", err)
	}

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	data, err := getData()
	if err != nil {
		log.Fatalf("Unable to get data: %v", err)
	}

	for _, address := range data.Addresses {
		msg := MessageDetails{
			SenderName:    data.From.Name,
			SenderAddress: data.From.Address,
			To:            address,
			Subject:       data.Subject,
			Body:          data.Body,
		}

		gmailMsg := createGmailMessage(msg)
		_, err = srv.Users.Messages.Send("me", &gmailMsg).Do()
		if err != nil {
			log.Fatalf("Unable to send email: %v", err)
		}

		log.Printf("Email sent to %s", msg.To)
	}

}
