package routerFuncs

import (
	"encoding/json"
	"net/http"
)

type RootMessage struct {
	Server string `json:"server"`
}

func Root(w http.ResponseWriter, r *http.Request) {
	message := RootMessage{
		Server: "Tahlilchi.uz",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}
