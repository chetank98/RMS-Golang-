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
