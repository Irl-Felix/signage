package util

import (
	"encoding/json"
	"io"
	"net/http"
)

// DecodeJSONBody decodes a JSON request body into the given destination struct.
// Returns an error if decoding fails or the body is empty.
func DecodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Body == nil {
		HandleError(w, nil, "Empty request body", http.StatusBadRequest)
		return io.EOF
	}
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		HandleError(w, err, "Invalid JSON input", http.StatusBadRequest)
		return err
	}
	return nil
}
