package main

import (
	"log"
	"net/http"

	"github.com/dineshkuncham/logger-service/data"
)

type requestPaylod struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {

	payload := requestPaylod{}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//validate user credentials

	err = app.Models.LogEntry.Insert(data.LogEntry{
		Name: payload.Name,
		Data: payload.Data,
	})
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: "Logged",
	}

	err = app.writeJSON(w, http.StatusAccepted, response)

	if err != nil {
		log.Printf("Error sending response: %s", err)
	}

}
