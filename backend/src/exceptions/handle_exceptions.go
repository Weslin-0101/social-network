package exceptions

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandleResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Fatal(err)
	}
}

func HandleError(w http.ResponseWriter, statusCode int, err error) {
	HandleResponse(w, statusCode, struct {
		Erro string `json:"error"`
	}{
		Erro: err.Error(),
	})
}