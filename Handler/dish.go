package Handler

import (
	"RMS/Database/DbHelper"
	"RMS/Middleware"
	"RMS/Models"
	"RMS/Utils"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func CreateDish(w http.ResponseWriter, r *http.Request) {
	restaurantID := chi.URLParam(r, "restaurantId")
	var data Models.DishRequest

	if parseErr := Utils.ParsreBody(r.Body, &data); parseErr != nil {
		Utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	exists, existsErr := DbHelper.IsDishExists(data.Name, restaurantID)
	if existsErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, existsErr, "failed to check dish existence")
		return
	}
	if exists {
		Utils.RespondError(w, http.StatusConflict, nil, "dish already exists")
		return
	}

	if saveErr := DbHelper.DishCreation(data, restaurantID); saveErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, saveErr, "failed to save dish")
		return
	}

	Utils.RespondJSON(w, http.StatusCreated, struct {
		Message string `json:"message"`
	}{"dish created successfully"})
}

func GetAllDishes(w http.ResponseWriter, r *http.Request) {
	dishes, getErr := DbHelper.GetAllDishes()

	if getErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get dishes")
		return
	}

	Utils.RespondJSON(w, http.StatusOK, dishes)
}

func GetAllDishesBySubAdmin(w http.ResponseWriter, r *http.Request) {
	userCtx := Middleware.UserContext(r)
	loggedUserID := userCtx.UserID

	dishes, getErr := DbHelper.GetAllDishesBySubAdmin(loggedUserID)
	if getErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get dishes")
		return
	}

	Utils.RespondJSON(w, http.StatusOK, dishes)
}

func DishesByRestaurant(w http.ResponseWriter, r *http.Request) {
	data := struct {
		RestaurantID string `json:"restaurantId" db:"restaurant_id" validate:"required"`
	}{}

	if parseErr := Utils.ParsreBody(r.Body, &data); parseErr != nil {
		Utils.RespondError(w, http.StatusBadRequest, parseErr, "failed to parse request body")
		return
	}

	dishes, getErr := DbHelper.DishesByRestaurant(data.RestaurantID)
	if getErr != nil {
		Utils.RespondError(w, http.StatusInternalServerError, getErr, "failed to get dishes")
		return
	}

	Utils.RespondJSON(w, http.StatusOK, dishes)
}
