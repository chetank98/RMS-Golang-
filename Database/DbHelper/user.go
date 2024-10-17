package DbHelper

import (
	"RMS/Database"
	"RMS/Models"
	"github.com/jmoiron/sqlx"
	"time"
)

func CreateUser(tx *sqlx.Tx, name, email, password, created_by string, role Models.Role) (string, error) {

	sqlQuery := `INSERT INTO users (name, email, password, created_by, role)
							VALUES (TRIM($1),TRIM($2),$3,$4,$5) returning id`

	var UserId string
	creErr := tx.Get(&UserId, sqlQuery, name, email, password, created_by, role)
	return UserId, creErr
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
