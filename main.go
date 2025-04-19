package main

import (
	"context"
	"encoding/base64"
	"log"
	"os"
	"strconv"
	"strings"

	_ "embed"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

//go:embed body.txt
var body string

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

	isDryRun, err := getIsDryRun()
	if err != nil {
		log.Fatalf("Unable to get dry run flag: %v", err)
	}

	for _, address := range data.Addresses {
		msg := MessageDetails{
			SenderName:    data.From.Name,
			SenderAddress: data.From.Address,
			To:            address,
			Subject:       data.Subject,
			Body:          strings.ReplaceAll(body, "\n", "<br>"),
		}

		gmailMsg := createGmailMessage(msg)

		if isDryRun {
			msgRaw, err := base64.URLEncoding.DecodeString(gmailMsg.Raw)
			if err != nil {
				log.Fatalf("Unable to decode email: %v", err)
			}

			log.Print("-----------------------------")
			log.Print(string(msgRaw))
			log.Print("-----------------------------")
			continue
		}

		_, err = srv.Users.Messages.Send("me", &gmailMsg).Do()
		if err != nil {
			log.Fatalf("Unable to send email: %v", err)
		}

		log.Printf("Email sent to %v", msg.To)

	}

}

func getIsDryRun() (bool, error) {
	isDryRun := false
	var err error

	if len(os.Args) > 1 {
		isDryRun, err = strconv.ParseBool(os.Args[1])
		if err != nil {
			return false, err
		}
	}

	return isDryRun, nil
}
