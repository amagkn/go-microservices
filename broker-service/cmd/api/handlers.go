package main

import (
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, _ *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeResponse(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readRequestBody(w, r, &requestPayload)
	if err != nil {
		_ = app.writeErrorResponse(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	default:
		_ = app.writeErrorResponse(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	var jsonFromService jsonResponse

	err := runRequest("POST", "http://authentication-service/authenticate", a, &jsonFromService)

	if err != nil {
		_ = app.writeErrorResponse(w, err)
		return
	}

	if jsonFromService.Error {
		_ = app.writeErrorResponse(w, err, http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    jsonFromService.Data,
	}

	_ = app.writeResponse(w, http.StatusAccepted, payload)
}
