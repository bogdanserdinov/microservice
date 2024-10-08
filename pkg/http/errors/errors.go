package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

func ServeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	var response struct {
		Error string `json:"error"`
	}

	response.Error = err.Error()

	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Println("failed to write json error response", err)
	}
}
