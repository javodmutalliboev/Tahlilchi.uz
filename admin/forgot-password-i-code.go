package admin

import (
	"log"
	"net/http"
	"strconv"

	"Tahlilchi.uz/response"
)

func ForgotPasswordICode(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		response.Res(w, "error", http.StatusBadRequest, err.Error())
		return
	}

	iCode := r.Form.Get("i-code")
	if iCode == "" {
		response.Res(w, "error", http.StatusBadRequest, "i-code not provided")
		return
	}

	authentication := iCodeAuth(r, iCode)
	if !authentication.status && authentication.message != "" {
		if authentication.message == "Forbidden" {
			response.Res(w, "error", http.StatusForbidden, authentication.message)
			return
		} else {
			response.Res(w, "error", http.StatusInternalServerError, authentication.message)
			return
		}
	}

	session, _ := Store.Get(r, "admin-forgot-password")
	session.Values["#i#-$code$-?authenticated?"] = true
	session.Save(r, w)

	response.Res(w, "success", http.StatusOK, "i-code authenticated")
}

func iCodeAuth(r *http.Request, iCode string) iCodeAuthRT {
	session, _ := Store.Get(r, "admin-forgot-password")

	iCodeI, err := strconv.Atoi(iCode)
	if err != nil {
		log.Println(err.Error())
		return iCodeAuthRT{status: false, message: err.Error()}
	}

	if siCode, ok := session.Values["i-code"].(int); siCode != iCodeI || !ok {
		return iCodeAuthRT{
			status:  false,
			message: "Forbidden",
		}
	}

	return iCodeAuthRT{status: true, message: ""}
}

type iCodeAuthRT struct {
	status  bool
	message string
}
