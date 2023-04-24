package main

import (
	"net/http"
)

func (app *Config) Broker(res http.ResponseWriter, req *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJson(res, http.StatusOK, payload)
}
