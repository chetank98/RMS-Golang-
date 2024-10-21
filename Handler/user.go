package Handler

import (
	"RMS/Database"
	"RMS/Database/DbHelper"
	"RMS/Middleware"
	"RMS/Models"
	"RMS/Utils"
	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("Welcome to the RMS"))

}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var data Models.UserRequest

	userCtx := Middleware.UserContext(r)
	createBy := userCtx.UserID
	role := Models.RoleUser

	if parseErr := Utils.ParsreBody(r.Body, &data); parseErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, parseErr, "Error in decoding the data")
		return
	}

	//validate := validator.New()
	//if err := validate.Struct(data); err != nil {
	//	Utils.RespondError(w, http.StatusBadRequest, err, "Validation failed")
	//	return
	//}

	exists, existErr := DbHelper.AlreadyUser(data.Email)
	if existErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, existErr, "Error in checking the user")
		return
	}

	if exists {
		Utils.RespondError(w, http.StatusConflict, nil, "User already exists")
		return
	}

	hashPassword, hashErr := Utils.HashPassword(data.Password)
	if hashErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, hashErr, "Failed to secure password")
	}

	if txErr := Database.Tx(func(tx *sqlx.Tx) error {
		userId, crtErr := DbHelper.CreateUser(tx, data.Name, data.Email, hashPassword, createBy, role)
		if crtErr != nil {
			return crtErr
		}
		return DbHelper.CreateUserAddress(tx, userId, data.Address)
	}); txErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, txErr, "failed to create user")
	}

	Utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{"user created successfully"})

}

func LoginUser(w http.ResponseWriter, r *http.Request) {

	var data Models.LoginRequest

	if parseErr := Utils.ParsreBody(r.Body, &data); parseErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, parseErr, "Error in parsing the request")
	}

	validate := validator.New()
	if err := validate.Struct(data); err != nil {
		Utils.RespondError(w, http.StatusBadRequest, err, "input validation failed")
		return
	}

	userId, role, getErr := DbHelper.GettingLoginDetails(data)
	if getErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to find user")
		return
	}

	if userId == "" || role == "" {
		Utils.RespondError(w, http.StatusBadRequest, nil, "user not found")
	}

	sessionID, startErr := DbHelper.SessionStart(userId)
	if startErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, startErr, "session not started")
		return
	}

	token, genErr := Utils.GenerateJWT(userId, sessionID, role)
	if genErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, genErr, "failed to generate token")
		return
	}

	Utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
		Token   string `json:"token"`
	}{"login successful", token})

}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	userCtx := Middleware.UserContext(r)
	sessionID := userCtx.SessionID

	if delErr := DbHelper.DeleteUserSession(sessionID); delErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, delErr, "failed to delete user session")
		return
	}

	Utils.RespondJSON(w, http.StatusOK, struct {
		Message string `json:"message"`
	}{"logout successful"})
}

func GetAllUsersByAdmin(w http.ResponseWriter, r *http.Request) {
	users, getErr := DbHelper.GetAllUsersByAdmin()

	if getErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get users")
		return
	}

	Utils.RespondJSON(w, http.StatusOK, users)
}

func GetAllUsersBySubAdmin(w http.ResponseWriter, r *http.Request) {
	userCtx := Middleware.UserContext(r)
	loggedUserID := userCtx.UserID

	users, getErr := DbHelper.GetAllUsersBySubAdmin(loggedUserID)
	if getErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get users")
		return
	}

	Utils.RespondJSON(w, http.StatusOK, users)
}
