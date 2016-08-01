package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/nlopes/slack"
	"github.com/subosito/twilio"
)

var (
	fromNumber = os.Getenv("FROM_NUMBER")
	toNumber   = os.Getenv("TO_NUMBER")

	slackToken     = os.Getenv("SLACK_TOKEN")
	channelToEnvoy = os.Getenv("SLACK_CHANNEL")

	twilioSid   = os.Getenv("TWILIO_SID")
	twilioToken = os.Getenv("TWILIO_TOKEN")
)

func main() {
	twilioClient := twilio.NewClient(twilioSid, twilioToken, nil)
	slackManager, err := NewSlackManager(slackToken, channelToEnvoy)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing Slack:", err)
		os.Exit(1)
	}

	go slackManager.Run(twilioClient)
	listenSMS(twilioClient, slackManager.api, slackManager.channelID, slackManager.userID)
}

func listenSMS(twilioClient *twilio.Client, slackClient *slack.Client, channelID, userID string) error {
	http.HandleFunc("/callback/sms", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing SMS callback")
		}

		msg := r.FormValue("Body")

		params := slack.NewPostMessageParameters()
		params.AsUser = true

		if _, _, err := slackClient.PostMessage(channelID, msg, params); err != nil {
			fmt.Fprintln(os.Stderr, "Error sending Slack message")
		}
	})

	fmt.Println("Listening for SMS callbacks on port 3000")
	return http.ListenAndServe(":3000", nil)
}

func sendSMS(client *twilio.Client, fromNumber, toNumber, msg string) error {
	params := twilio.MessageParams{Body: msg}
	if _, _, err := client.Messages.Send(fromNumber, toNumber, params); err != nil {
		return err
	}

	return nil
}
