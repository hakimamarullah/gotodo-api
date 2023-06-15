package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/gotodo/models"
)

func ResponseJSON(w http.ResponseWriter, payload models.ResponseBody) {
	if payload.Message == "" {
		payload.Message = "Success"
	}

	if payload.Code == 0 {
		payload.Code = http.StatusOK
	}

	response, _ := json.Marshal(payload)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(payload.Code)
	w.Write(response)
}
