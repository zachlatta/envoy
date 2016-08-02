package main

import (
	"fmt"
	"os"
)

var (
	port = os.Getenv("PORT")

	fromNumber = os.Getenv("FROM_NUMBER")
	toNumber   = os.Getenv("TO_NUMBER")

	slackToken     = os.Getenv("SLACK_TOKEN")
	channelToEnvoy = os.Getenv("SLACK_CHANNEL")

	twilioSid   = os.Getenv("TWILIO_SID")
	twilioToken = os.Getenv("TWILIO_TOKEN")
)

func main() {
	// Set port to 3000 by default
	if port == "" {
		port = "3000"
	}

	msgService := NewTwilioService(twilioSid, twilioToken, fromNumber, toNumber)
	slackManager, err := NewSlackManager(slackToken, channelToEnvoy)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing Slack:", err)
		os.Exit(1)
	}

	go slackManager.Run(msgService)
	msgService.Listen(":"+port, "/callback/sms", func(msg string) {
		if err := slackManager.SendMessage(msg); err != nil {
			fmt.Fprintln(os.Stderr, "Error sending Slack message:", err)
		}
	})
}
