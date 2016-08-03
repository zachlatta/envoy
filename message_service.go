package main

import (
	"fmt"
	"github.com/subosito/twilio"
	"net/http"
	"os"
)

type MessageService interface {
	SendMessage(msg string) error
	Listen(port, path string, callback func(msg string)) error
}

type TwilioService struct {
	api *twilio.Client

	fromNumber string
	toNumber   string
}

func NewTwilioService(sid, token, fromNumber, toNumber string) *TwilioService {
	return &TwilioService{api: twilio.NewClient(sid, token, nil)}
}

func (s *TwilioService) SendMessage(msg string) error {
	params := twilio.MessageParams{Body: msg}
	if _, _, err := s.api.Messages.Send(s.fromNumber, s.toNumber, params); err != nil {
		return err
	}

	return nil
}

func (s *TwilioService) Listen(port, path string, callback func(msg string)) error {
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			fmt.Fprintln(os.Stderr, "Error parsing SMS callback")
		}

		msg := r.FormValue("Body")

		callback(msg)

		w.Write(nil)
	})

	return http.ListenAndServe(port, nil)
}
