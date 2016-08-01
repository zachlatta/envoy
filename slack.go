package main

import (
	"errors"

	"fmt"
	"github.com/nlopes/slack"
	"os"
)

var (
	ErrChannelNotFound = errors.New("channel not found")
)

type SlackManager struct {
	api *slack.Client

	userID    string
	channelID string
}

func NewSlackManager(token, channelName string) (*SlackManager, error) {
	api := slack.New(token)

	userID, err := currentUserID(api)
	if err != nil {
		return nil, err
	}

	channelID, err := channelID(api, channelName)
	if err != nil {
		return nil, err
	}

	return &SlackManager{
		api:       api,
		userID:    userID,
		channelID: channelID,
	}, nil
}

func (m *SlackManager) Run(msgService MessageService) {
	rtm := m.api.NewRTM()
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

				if msg.Channel != m.channelID || msg.User == m.userID {
					continue
				}

				if err := msgService.SendMessage(msg.Text); err != nil {
					fmt.Fprintln(os.Stderr, "Failed to send SMS:", msg.Text)
				}

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				break loop
			}
		}
	}
}

func (m *SlackManager) SendMessage(msg string) error {
	params := slack.NewPostMessageParameters()
	params.AsUser = true

	if _, _, err := m.api.PostMessage(m.channelID, msg, params); err != nil {
		return err
	}

	return nil
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

	return "", ErrChannelNotFound
}
