package DbHelper

import (
	"RMS/Database"
	"RMS/Models"
)

func IsRestaurantExists(name, address string) (bool, error) {
	SQL := ` SELECT count(id) > 0 as is_exist
				FROM restaurants
				WHERE name = TRIM($1)
				  AND address = TRIM($2)
				  AND archived_at IS NULL`

	var check bool
	checkErr := Database.DBConnection.Get(&check, SQL, name, address)
	return check, checkErr
}

func CreateRestaurants(data Models.RestaurantsRequest, userId string) error {
	sqlQuery := `INSERT INTO restaurants(name, address, latitude, longitude, created_by) 
							VALUES (TRIM($1),TRIM($2),$3,$4,$5) `

	_, creErr := Database.DBConnection.Exec(sqlQuery, data.Name, data.Address, data.Latitude, data.Longitude, userId)
	return creErr

}

func GetallRestaurants() ([]Models.Restaurant, error) {

	sqlQuery := `SELECT id, name, address, latitude, longitude, created_by
    					FROM restaurants 
						WHERE archived_at IS NULL`

	restaurant := make([]Models.Restaurant, 0)
	getErr := Database.DBConnection.Select(&restaurant, sqlQuery)
	return restaurant, getErr

}

func GetAllRestaurantsBySubAdmin(loggedUserID string) ([]Models.Restaurant, error) {

	sqlQuery := `SELECT id, name, address, latitude, longitude, created_by
    					FROM restaurants 
						WHERE created_by = $1
						AND archived_at IS NULL`

	restaurants := make([]Models.Restaurant, 0)
	fetchErr := Database.DBConnection.Select(&restaurants, sqlQuery, loggedUserID)
	return restaurants, fetchErr

}
