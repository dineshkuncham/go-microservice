package main

import (
	"fmt"
	"log"
	"net/http"
)

type requestPaylod struct {
	From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func (app *Config) SendMail(w http.ResponseWriter, r *http.Request) {
	payload := requestPaylod{}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}

	msg := Message{
		From:    payload.From,
		To:      payload.To,
		Subject: payload.Subject,
		Data:    payload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		log.Println(err)
		app.errorJSON(w, err)
		return
	}
	response := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Sent to %s", payload.To),
	}

	err = app.writeJSON(w, http.StatusAccepted, response)

	if err != nil {
		log.Panic(err)
	}

}
