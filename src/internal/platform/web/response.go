package web

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func Response(w http.ResponseWriter, data interface{}, statusCode int) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.Wrapf(err, "marshalling data: %v", data)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(jsonData); err != nil {
		return errors.Wrapf(err, "writing data: %v", data)
	}
	return nil
}

// RespondError ResponseError know how to handle errors going out to the client.
func RespondError(w http.ResponseWriter, err error) error {
	if errWeb, ok := err.(*Error); ok {
		resp := ErrorResponse{
			Error: errWeb.Err.Error(),
		}
		return Response(w, resp, errWeb.Code)
	}

	resp := ErrorResponse{
		Error: http.StatusText(http.StatusInternalServerError),
	}
	return Response(w, resp, http.StatusInternalServerError)
}
