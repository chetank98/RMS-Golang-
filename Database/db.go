package Database

import (
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
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

func Tx(fn func(tx *sqlx.Tx) error) error {
	tx, err := DBConnection.Beginx()
	if err != nil {
		return fmt.Errorf("failed to start a transaction: %+v", err)
	}
	defer func() {
		if err != nil {
			if rollBackErr := tx.Rollback(); rollBackErr != nil {
				logrus.Errorf("failed to rollback tx: %s", rollBackErr)
			}
			return
		}
		if commitErr := tx.Commit(); commitErr != nil {
			logrus.Errorf("failed to commit tx: %s", commitErr)
		}
	}()
	err = fn(tx)
	return err
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

func ShutdownDatabase() error {
	return DBConnection.Close()
}
