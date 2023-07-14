package main

import "net/http"

func (app *Config) SendMail(res http.ResponseWriter, req *http.Request) {
	type mailMessage struct {
		From    string `json:"from"`
		To      string `json:"to"`
		Subject string `json:"subject"`
		Message string `json:message`
	}

	var requestPayload mailMessage

	err := app.readJson(res, req, &requestPayload)

	if err != nil {
		app.errorJson(res, err)
		return
	}

	msg := Message{
		From:    requestPayload.From,
		To:      requestPayload.To,
		Subject: requestPayload.Subject,
		Data:    requestPayload.Message,
	}

	err = app.Mailer.SendSMTPMessage(msg)
	if err != nil {
		app.errorJson(res, err)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "send to " + requestPayload.To,
	}

	app.writeJson(res, http.StatusOK, payload)
}
