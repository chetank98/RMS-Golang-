package Handler

import (
	"RMS/Database/DbHelper"
	"RMS/Middleware"
	"RMS/Models"
	"RMS/Utils"
	"net/http"
)

func CreateRestaurants(w http.ResponseWriter, r *http.Request) {
	var data Models.RestaurantsRequest

	userCtx := Middleware.UserContext(r)
	createBy := userCtx.UserID

	if parseErr := Utils.ParsreBody(r.Body, &data); parseErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, parseErr, "Error in decoding the data")
		return
	}

	exists, existErr := DbHelper.IsRestaurantExists(data.Name, createBy)
	if existErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, existErr, "Error in checking the restaurants")
		return
	}

	if exists {
		Utils.RespondError(w, http.StatusConflict, nil, "restaurant already exists")
		return
	}

	if saveErr := DbHelper.CreateRestaurants(data, createBy); saveErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to save restaurant")
		return
	}

	Utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{"restaurant created successfully"})

}

func GetallRestaurants(w http.ResponseWriter, r *http.Request) {

	restaurants, getErr := DbHelper.GetallRestaurants()

	if getErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get restaurants")
		return
	}

	Utils.RespondJSON(w, http.StatusOK, restaurants)

}

func GetAllRestaurantsBySubAdmin(w http.ResponseWriter, r *http.Request) {

	userCtx := Middleware.UserContext(r)
	loggedUserID := userCtx.UserID

	restaurants, getErr := DbHelper.GetAllRestaurantsBySubAdmin(loggedUserID)
	if getErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get restaurants")
		return
	}

	Utils.RespondJSON(w, http.StatusOK, restaurants)

}
