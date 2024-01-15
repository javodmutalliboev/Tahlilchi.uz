package admin

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"Tahlilchi.uz/response"
)

func ForgotPasswordICode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := r.ParseForm()
	if err != nil {
		res := response.Response{
			Status:     "error",
			StatusCode: http.StatusBadRequest,
			Data:       err.Error(),
		}

		json.NewEncoder(w).Encode(res)
		return
	}

	iCode := r.Form.Get("i-code")
	if iCode == "" {
		res := response.Response{
			Status:     "error",
			StatusCode: http.StatusBadRequest,
			Data:       "i-code not provided",
		}

		json.NewEncoder(w).Encode(res)
	}

	authentication := iCodeAuth(r, iCode)
	if !authentication.status && authentication.message != "" {
		if authentication.message == "Forbidden" {
			res := response.Response{
				Status:     "error",
				StatusCode: http.StatusForbidden,
				Data:       authentication.message,
			}

			json.NewEncoder(w).Encode(res)
			return
		} else {
			res := response.Response{
				Status:     "error",
				StatusCode: http.StatusInternalServerError,
				Data:       authentication.message,
			}

			json.NewEncoder(w).Encode(res)
			return
		}
	}

	session, _ := Store.Get(r, "admin-forgot-password")
	session.Values["#i#-$code$-?authenticated?"] = true
	session.Save(r, w)

	res := response.Response{
		Status:     "success",
		StatusCode: http.StatusOK,
		Data:       "i-code authenticated",
	}

	json.NewEncoder(w).Encode(res)
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
