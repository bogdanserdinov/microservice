package errors

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Error string `json:"error"`
}

func ServeError(w http.ResponseWriter, status int, err error) {
	w.WriteHeader(status)

	response := Response{
		Error: err.Error(),
	}

	if err = json.NewEncoder(w).Encode(response); err != nil {
		log.Println("failed to write json error response", err)
	}
}
