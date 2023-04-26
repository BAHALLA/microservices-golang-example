package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (app *Config) readJson(res http.ResponseWriter, req *http.Request, data any) error {
	maxByte := 1048576

	req.Body = http.MaxBytesReader(res, req.Body, int64(maxByte))

	dec := json.NewDecoder(req.Body)
	err := dec.Decode(data)

	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must be only a single JSON value")
	}
	return nil

}

func (app *Config) writeJson(res http.ResponseWriter, status int, data any, headers ...http.Header) error {

	out, err := json.Marshal(data)

	if err != nil {
		return err
	}

	if len(headers) > 0 {

		for key, value := range headers[0] {
			res.Header()[key] = value
		}
	}
	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(status)
	_, err = res.Write(out)
	if err != nil {
		return err
	}

	return nil

}

func (app *Config) errorJson(res http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	var payload jsonResponse
	payload.Error = true
	payload.Message = err.Error()

	return app.writeJson(res, statusCode, payload)
}
