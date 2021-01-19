package web

import (
	// Core packages
	"encoding/json"
	"net/http"

	// Third party packages
	"github.com/pkg/errors"
)

// Respond converts a Go value to JSON and sends it to the client.
func Respond(w http.ResponseWriter, data interface{}, statusCode int) error {

	// In the case of Updates there will be no content to respond with
	if statusCode == http.StatusNoContent {
		w.WriteHeader(statusCode)
		return nil
	}

	// Convert the response value to JSON.
	// https://golang.org/pkg/encoding/json/#Marshal
	res, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Respond with the provided JSON.
	// https://golang.org/pkg/net/http/#Request.Write
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(statusCode)
	if _, err := w.Write(res); err != nil {
		return err
	}

	return nil
}

// RespondError sends an error reponse back to the client.
func RespondError(w http.ResponseWriter, err error) error {

	// If the error was of the type *Error, the handler has
	// a specific status code and error to return.
	if webErr, ok := errors.Cause(err).(*Error); ok {
		er := ErrorResponse{
			Error:  webErr.Err.Error(),
			Fields: webErr.Fields,
		}
		if err := Respond(w, er, webErr.Status); err != nil {
			return err
		}
		return nil
	}

	// If not, the handler sent any arbitrary error value so use 500.
	er := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	if err := Respond(w, er, http.StatusInternalServerError); err != nil {
		return err
	}
	return nil
}
