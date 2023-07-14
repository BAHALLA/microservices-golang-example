package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(res http.ResponseWriter, req *http.Request) {
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(res, req, &requestPayload)

	if err != nil {
		app.errorJson(res, err, http.StatusBadRequest)
		return
	}

	user, err := app.models.User.GetByEmail(requestPayload.Email)

	if err != nil {
		app.errorJson(res, errors.New("invalide credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)

	if err != nil || !valid {
		app.errorJson(res, errors.New("invalide credentials"), http.StatusUnauthorized)
		return
	}

	err = app.logRequest("authentication", fmt.Sprintf("%s is logged", user.Email))

	if err != nil {
		app.errorJson(res, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged as user %s", requestPayload.Email),
		Data:    user,
	}
	app.writeJson(res, http.StatusAccepted, payload)

}

func (app *Config) logRequest(name, data string) error {

	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.MarshalIndent(entry, "", "\t")

	request, err := http.NewRequest("POST", "http://logger-service:8084/log", bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	client := &http.Client{}
	_, err = client.Do(request)

	if err != nil {
		return err
	}

	return nil
}
