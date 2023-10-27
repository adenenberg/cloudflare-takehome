package handlers

import (
	"encoding/json"
	"net/http"
)

type data struct {
	ShortenedURL interface{} `json:"data,omitempty"`
	Statuscode   int         `json:"status"`
	Message      string      `json:"message"`
}

func SuccessResponse(payload interface{}, w http.ResponseWriter) {
	_, err := json.Marshal(payload)
	resp := &data{ShortenedURL: payload, Statuscode: 200, Message: "success"}
	if err != nil {
		// err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
