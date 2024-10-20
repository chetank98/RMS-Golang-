package Handler

import (
	"RMS/Database/DbHelper"
	"RMS/Middleware"
	"RMS/Models"
	"RMS/Utils"
	"net/http"
)

func SubAdminCreation(w http.ResponseWriter, r *http.Request) {

	var data Models.SubAdminRequest

	userCtx := Middleware.UserContext(r)
	createdBy := userCtx.UserID
	role := Models.RoleSubAdmin

	if parseErr := Utils.ParsreBody(r.Body, &data); parseErr != nil {
		Utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	exists, existsErr := DbHelper.AlreadyUser(data.Email)
	if existsErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check user existence")
		return
	}
	if exists {
		Utils.RespondError(w, http.StatusConflict, nil, "sub-admin already exists")
		return
	}

	hashedPassword, hasErr := Utils.HashPassword(data.Password)
	if hasErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, hasErr, "failed to secure password")
		return
	}

	if saveErr := DbHelper.SubAdminCreation(data.Name, data.Email, hashedPassword, createdBy, role); saveErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to create sub-admin")
		return
	}

	Utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{"sub-admin created successfully"})
}

func GetAllSubAdmins(w http.ResponseWriter, r *http.Request) {
	subAdmins, getErr := DbHelper.GetAllSubAdmins()

	if getErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get sub-admin")
		return
	}

	Utils.RespondJSON(w, http.StatusOK, subAdmins)
}
