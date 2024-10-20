package DbHelper

import (
	"RMS/Database"
	"RMS/Models"
)

func SubAdminCreation(name, email, password, createdBy string, role Models.Role) error {

	sqlQuery := `INSERT INTO users (name, email, password, created_by, role)
			  VALUES (TRIM($1), TRIM($2), $3, $4, $5) RETURNING id`

	var userID string
	crtErr := Database.DBConnection.Get(&userID, sqlQuery, name, email, password, createdBy, role)
	return crtErr

}

func GetAllSubAdmins() ([]Models.SubAdmin, error) {
	SQL := `SELECT id,
				   name,
				   email,
				   role,
				   created_by
			FROM users
				WHERE role = 'sub-admin' 
				AND archived_at IS NULL`

	subAdmins := make([]Models.SubAdmin, 0)
	fetchErr := Database.DBConnection.Select(&subAdmins, SQL)
	return subAdmins, fetchErr
}
