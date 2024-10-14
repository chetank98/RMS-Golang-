package Database

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var (
	DBConnection *sqlx.DB
)

type SSLMode string

const (
	SSLModeDisable SSLMode = "disable"
)

func ConnectAndMigrate(host, port, databasename, user, password string, sslMode SSLMode) error {

	connection := fmt.Sprintf("host=%s port=%s user=%s password=%s  dbname=%s sslmode=%s",
		host, port, user, password, databasename, sslMode)

	DB, err := sqlx.Open("postgres", connection)
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}
	DBConnection = DB
	return migrateUp(DB)

}

func ShutdownDatabase() error {
	return DBConnection.Close()
}

func migrateUp(db *sqlx.DB) error {
	// migrate the database and handle the migration logic
	driver, driErr := postgres.WithInstance(db.DB, &postgres.Config{})
	if driErr != nil {
		return driErr
	}
	m, migErr := migrate.NewWithDatabaseInstance(
		"file://Database/Migrations/",
		"postgres", driver)
	if migErr != nil {
		return migErr
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
