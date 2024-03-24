package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

type requestPaylod struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	payload := requestPaylod{}

	err := app.readJSON(w, r, &payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	//validate user credentials

	user, err := app.Repo.GetByEmail(payload.Email)
	if err != nil {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}
	valid, err := app.Repo.PasswordMatches(payload.Password, *user)
	if err != nil || !valid {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	//log request
	err = app.logRequest("authentication-service", fmt.Sprintf("%s logged in", user.Email))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	response := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}

	err = app.writeJSON(w, http.StatusAccepted, response)

	if err != nil {
		log.Printf("Error sending response: %s", err)
	}

}

func (app *Config) logRequest(name, data string) error {
	entry := struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}{
		Name: name,
		Data: data,
	}
	// create some json we'll send to the auth microservice
	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	// call the service
	request, err := http.NewRequest("POST", "http://logger-service/log", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	_, err = app.Client.Do(request)
	if err != nil {
		return err
	}
	return nil
}
