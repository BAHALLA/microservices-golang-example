package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPaylaod struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(res http.ResponseWriter, req *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJson(res, http.StatusOK, payload)
}

func (app *Config) HandleRequest(res http.ResponseWriter, req *http.Request) {
	var requestPayload RequestPaylaod

	err := app.readJson(res, req, &requestPayload)

	if err != nil {
		app.errorJson(res, err)
		return
	}

	switch requestPayload.Action {

	case "auth":
		app.authenticate(res, requestPayload.Auth)
	case "log":
		app.logItem(res, requestPayload.Log)
	default:
		app.errorJson(res, errors.New("unknown action"))
	}

}

func (app *Config) authenticate(res http.ResponseWriter, authPayload AuthPayload) {

	jsonData, _ := json.MarshalIndent(authPayload, "", "\t")

	request, err := http.NewRequest("POST", "http://authentification-service:8083/authenticate", bytes.NewBuffer(jsonData))

	if err != nil {
		app.errorJson(res, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJson(res, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {

		app.errorJson(res, errors.New("invalide credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJson(res, errors.New("error calling auth service"))
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)

	if err != nil {
		app.errorJson(res, err)
		return
	}

	if jsonFromService.Error {
		app.errorJson(res, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJson(res, http.StatusAccepted, payload)

}

func (app *Config) logItem(res http.ResponseWriter, logPLogPayload LogPayload) {
	jsonData, _ := json.MarshalIndent(logPLogPayload, "", "\t")

	request, err := http.NewRequest("POST", "http://logger-service:8084/log", bytes.NewBuffer(jsonData))

	if err != nil {
		app.errorJson(res, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		app.errorJson(res, err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		app.errorJson(res, err)
		return
	}

	var jsonResponse jsonResponse

	jsonResponse.Error = false
	jsonResponse.Message= "logged"

	app.writeJson(res, http.StatusOK, jsonResponse)

}
