package main

import (
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

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged as user %s", requestPayload.Email),
		Data:    user,
	}
	app.writeJson(res, http.StatusAccepted, payload)

}
