package DbHelper

import (
	"RMS/Database"
	"RMS/Models"
	"RMS/Utils"
	"github.com/jmoiron/sqlx"
	"time"
)

func AlreadyUser(email string) (bool, error) {
	sqlQuery := `SELECT count(id)>0 as exists FROM users
					WHERE email = TRIM($1)
					AND archived_at IS NULL `
	var check bool
	existErr := Database.DBConnection.Get(&check, sqlQuery, email)
	return check, existErr
}

func CreateUser(tx *sqlx.Tx, name, email, password, createdBy string, role Models.Role) (string, error) {

	sqlQuery := `INSERT INTO users (name, email, password, created_by, role)
							VALUES (TRIM($1),TRIM($2),$3,$4,$5) returning id`

	var userId string
	creErr := tx.Get(&userId, sqlQuery, name, email, password, createdBy, role)
	return userId, creErr
}

func CreateUserAddress(tx *sqlx.Tx, userId string, addresses []Models.UserAddress) error {

	sqlQuery := `INSERT INTO address(user_id,address, latitude, longitude) VALUES `

	values := make([]interface{}, 0)
	for i := range addresses {
		values = append(values,
			userId,
			addresses[i].Address,
			addresses[i].Latitude,
			addresses[i].Longitude,
		)
	}

	sqlQuery = Utils.SetupBindVars(sqlQuery, "(?,?,?,?)", len(addresses))
	_, err := tx.Exec(sqlQuery, values...)
	return err

}

func SessionStart(userId string) (string, error) {
	sqlQuery := `INSERT into user_sessions(user_id)
					VALUES ($1) returning session_id`
	var sessionId string
	createErr := Database.DBConnection.Get(&sessionId, sqlQuery, userId)
	return sessionId, createErr
}

func DeleteUserSession(sessionID string) error {
	sqlQuery := `UPDATE user_sessions
			  SET archived_at = NOW()
			  WHERE user_id = $1
			    AND archived_at IS NULL`

	_, delErr := Database.DBConnection.Exec(sqlQuery, sessionID)
	return delErr
}

func GettingLoginDetails(data Models.LoginRequest) (string, Models.Role, error) {
	sqlQuery := `SELECT id, password, role 
					FROM users 
					WHERE email = TRIM($1)
						AND archived_at IS NULL`
	var body Models.LoginDetail
	if getErr := Database.DBConnection.Get(&body, sqlQuery, data.Email); getErr != nil {
		return "", "", getErr
	}
	if passwordErr := Utils.CheckPassword(data.Password, body.Password); passwordErr != nil {
		return "", "", passwordErr
	}
	return body.ID, body.Role, nil
}

func GetArchivedAt(sessionId string) (*time.Time, error) {

	var archivedAt *time.Time
	sqlQuery := `SELECT archived_at 
              FROM user_sessions 
              WHERE session_id = $1
              	AND archived_at IS NULL`

	getErr := Database.DBConnection.Get(&archivedAt, sqlQuery, sessionId)
	return archivedAt, getErr

}

func GetAllUsersByAdmin() ([]Models.User, error) {
	sqlQuery := `SELECT id, name, email, role 
			FROM users
    	      WHERE role = 'User' 
    	        AND archived_at IS NULL`

	users := make([]Models.User, 0)
	if fetchErr := Database.DBConnection.Select(&users, sqlQuery); fetchErr != nil {
		return users, fetchErr
	}

	sqlQuery = `SELECT id, address, latitude, longitude, user_id 
			FROM address
    	      WHERE archived_at IS NULL`

	addresses := make([]Models.Address, 0)
	if fetchErr := Database.DBConnection.Select(&addresses, sqlQuery); fetchErr != nil {
		return users, fetchErr
	}

	addressMap := make(map[string][]Models.Address)
	for _, addr := range addresses {
		addressMap[addr.UserID] = append(addressMap[addr.UserID], addr)
	}

	for i := range users {
		if userAddresses, exists := addressMap[users[i].ID]; exists {
			users[i].Address = userAddresses
		}
	}

	return users, nil
}

func GetAllUsersBySubAdmin(loggedUserID string) ([]Models.User, error) {
	sqlQuery := `SELECT id, name, email, role 
			FROM users
    	      WHERE created_by = $1
    	        AND archived_at IS NULL`

	users := make([]Models.User, 0)
	if fetchErr := Database.DBConnection.Select(&users, sqlQuery, loggedUserID); fetchErr != nil {
		return users, fetchErr
	}

	sqlQuery = `SELECT a.id, a.address, a.latitude, a.longitude, a.user_id
			FROM address a
					 JOIN users u on a.user_id = u.id
			WHERE created_by = $1
			  AND a.archived_at IS NULL
			  AND u.archived_at IS NULL`

	addresses := make([]Models.Address, 0)
	if fetchErr := Database.DBConnection.Select(&addresses, sqlQuery, loggedUserID); fetchErr != nil {
		return users, fetchErr
	}

	addressMap := make(map[string][]Models.Address)
	for _, addr := range addresses {
		addressMap[addr.UserID] = append(addressMap[addr.UserID], addr)
	}

	for i := range users {
		if userAddresses, exists := addressMap[users[i].ID]; exists {
			users[i].Address = userAddresses
		}
	}

	return users, nil
}

func GetUserCoordinates(userAddressID string) (Models.Coordinates, error) {
	SQL := `SELECT latitude, longitude 
              FROM address 
              WHERE id = $1
              	AND archived_at IS NULL`

	var coordinates Models.Coordinates
	getErr := Database.DBConnection.Get(&coordinates, SQL, userAddressID)
	return coordinates, getErr
}

func GetRestaurantCoordinates(restaurantAddressID string) (Models.Coordinates, error) {
	SQL := `SELECT latitude, longitude 
              FROM restaurants 
              WHERE id = $1
              	AND archived_at IS NULL`

	var coordinates Models.Coordinates
	getErr := Database.DBConnection.Get(&coordinates, SQL, restaurantAddressID)
	return coordinates, getErr
}

func CalculateDistance(userCoordinates, restaurantCoordinates Models.Coordinates) (float64, error) {
	args := []interface{}{userCoordinates.Latitude, userCoordinates.Longitude,
		restaurantCoordinates.Latitude, restaurantCoordinates.Longitude}

	SQL := `SELECT ROUND(
						   (earth_distance(
									ll_to_earth($1, $2),
									ll_to_earth($3, $4)
							) / 1000.0)::numeric, 1
				   ) AS distance_km`

	var distance float64
	getErr := Database.DBConnection.Get(&distance, SQL, args...)
	return distance, getErr
}
