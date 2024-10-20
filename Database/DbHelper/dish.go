package DbHelper

import (
	"RMS/Database"
	"RMS/Models"
)

func IsDishExists(name, restaurantID string) (bool, error) {
	SQL := `SELECT count(id) > 0 as is_exist
				FROM dishes
				WHERE name = TRIM($1)
				  AND restaurant_id = $2
				  AND archived_at IS NULL`

	var check bool
	checkErr := Database.DBConnection.Get(&check, SQL, name, restaurantID)
	return check, checkErr
}

func DishCreation(body Models.DishCreation, restaurantID string) error {
	SQL := `INSERT INTO dishes (name, price, restaurant_id)
				VALUES (TRIM($1), $2, $3)`

	_, createErr := Database.DBConnection.Exec(SQL, body.Name, body.Price, restaurantID)
	return createErr
}

func GetAllDishes() ([]Models.DishCreation, error) {
	SQL := `SELECT id, name, price, restaurant_id
				FROM dishes
				WHERE archived_at IS NULL`

	dishes := make([]Models.DishCreation, 0)
	FetchErr := Database.DBConnection.Select(&dishes, SQL)
	return dishes, FetchErr
}

func GetAllDishesBySubAdmin(loggedUserID string) ([]Models.DishCreation, error) {
	SQL := `SELECT d.id, d.name, d.price, d.restaurant_id
				FROM dishes d
						 INNER JOIN restaurants r on d.restaurant_id = r.id
				WHERE d.archived_at IS NULL
				  AND r.created_by = $1`

	dishes := make([]Models.DishCreation, 0)
	fetchErr := Database.DBConnection.Select(&dishes, SQL, loggedUserID)
	return dishes, fetchErr
}

func DishesByRestaurant(restaurantID string) ([]Models.DishCreation, error) {
	SQL := `SELECT id, name, price, restaurant_id
				FROM dishes
				WHERE restaurant_id = $1
				  AND archived_at IS NULL`

	dishes := make([]Models.DishCreation, 0)
	fetchErr := Database.DBConnection.Select(&dishes, SQL, restaurantID)
	return dishes, fetchErr
}
