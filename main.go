package main

import (
	"errors"
	"fmt"
	"log"
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
	slackClient := slack.New(slackToken)
	logger := log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)
	slack.SetLogger(logger)

	twilioClient := twilio.NewClient(twilioSid, twilioToken, nil)

	userID, err := currentUserID(slackClient)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error retrieving authenticated user's information:", err)
		os.Exit(1)
	}

	channelID, err := channelID(slackClient, channelToEnvoy)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error retrieving given channel's information:", err)
		os.Exit(1)
	}

	go runSlack(slackClient, twilioClient, channelID, userID)
	listenSMS(twilioClient, slackClient, channelID, userID)
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

func runSlack(slackClient *slack.Client, twilioClient *twilio.Client, channelID, userID string) {
	rtm := slackClient.NewRTM()
	go rtm.ManageConnection()

loop:
	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.ConnectedEvent:
				fmt.Println("Connected to Slack! Connection count:", ev.ConnectionCount)

			case *slack.MessageEvent:
				msg := msg.Data.(*slack.MessageEvent).Msg

				if msg.Channel != channelID || msg.User == userID {
					continue
				}

				if err := sendSMS(twilioClient, fromNumber, toNumber, msg.Text); err != nil {
					fmt.Fprintln(os.Stderr, "Failed to send SMS:", msg.Text)
				}

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break loop
			}
		}
	}
}

func currentUserID(client *slack.Client) (string, error) {
	resp, err := client.AuthTest()
	if err != nil {
		return "", err
	}

	return resp.UserID, nil
}

func channelID(client *slack.Client, channelName string) (string, error) {
	channels, err := client.GetChannels(true)
	if err != nil {
		return "", err
	}

	for _, channel := range channels {
		if channel.Name == channelName {
			return channel.ID, nil
		}
	}

	return "", errors.New("channel not found")
}

func sendSMS(client *twilio.Client, fromNumber, toNumber, msg string) error {
	params := twilio.MessageParams{Body: msg}
	if _, _, err := client.Messages.Send(fromNumber, toNumber, params); err != nil {
		return err
	}

	return nil
}
