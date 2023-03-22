package app

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/pressly/goose"
	"github.com/spf13/viper"
	"os"
	"path"
)

var (
	db *sql.DB
)

func connectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("db.host"),
		"5432", // For localhost and docker inner network
		viper.GetString("db.user"),
		viper.GetString("db.pass"),
		viper.GetString("db.name"))
}

// InitDatabase For Vanilla SQL
func InitDatabase() (*sql.DB, error) {

	enabled := viper.GetBool("db.enabled")
	if !enabled {
		return nil, nil
	}

	p := path.Join(".", "migrations")
	if _, err := os.Stat(p); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(p, os.ModePerm)
		if err != nil {
			return nil, errors.Wrap(err, "cannot mkdir for migrations")
		}
	}

	driver := viper.GetString("db.driver")
	localDB, err := sql.Open(driver, connectionString())
	if err != nil {
		return nil, errors.Wrap(err, "cannot open database connection")
	}

	db = localDB

	// Test
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrapf(err, "Cannot ping to database: %s", connectionString())
	}
	log.Info("Successfully ping to database")

	log.Info("Database connection was created: %s", connectionString())
	return localDB, nil
}

func RunMigrations(rootDir ...string) error {

	enabled := viper.GetBool("db.enabled")
	if !enabled {
		return nil
	}

	basePath := "."
	if len(rootDir) != 0 {
		basePath = rootDir[0]
	}

	err := goose.Up(db, basePath+"/migrations")
	if err != nil {
		return err
	}
	return nil
}
