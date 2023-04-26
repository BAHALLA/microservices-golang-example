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
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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