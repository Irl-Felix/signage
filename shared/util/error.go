package util

import (
	"log"
	"net/http"
)

// HandleError standardizes error logging and HTTP response.
func HandleError(w http.ResponseWriter, err error, message string, status int) {
	if err != nil {
		log.Printf("[ERROR] %s: %v", message, err)
	} else {
		log.Printf("[ERROR] %s", message)
	}
	http.Error(w, message, status)
}

// Optionally, for non-HTTP contexts:
func LogError(err error, message string) {
	if err != nil {
		log.Printf("[ERROR] %s: %v", message, err)
	} else {
		log.Printf("[ERROR] %s", message)
	}
}
