package Handler

import (
	"RMS/Middleware"
	"RMS/Models"
	"RMS/Services"
	"RMS/Utils"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Welcome to the RMS"))

}

func CreateUser(w http.ResponseWriter, r http.Request) {

	var data Models.User

	userCtx := Middleware.UserContext(r)
	createBy := userCtx.UserID
	role := Models.RoleUser

	if parseErr := Utils.ParsreBody(r.Body, &data); parseErr != nil {
		return Utils.RespondError(w, http.StatusInternalServerError, parseErr, "Error in decoding the data")
	}

	err := Services.CreateUser()

}
