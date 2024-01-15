package routerFuncs

import (
	"net/http"

	"Tahlilchi.uz/response"
)

type RootMessage struct {
	Server string `json:"server"`
}

func Root(w http.ResponseWriter, r *http.Request) {

	response.Res(w, "success", http.StatusOK, "Tahlilchi.uz server")

}
