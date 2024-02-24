package helpers

import (
	"net/http"
)

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := JsonEnvelope{"error": message}

	err := WriteJSON(w, status, env, nil)
	if err != nil {
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := "the server encountered a problem and could not process your request"
	ErrorResponse(w, r, http.StatusInternalServerError, message)
}
