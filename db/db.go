package db

import (
	"database/sql"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func SetupConnection() (*sql.DB, error) {
	dbType := viper.GetString("db.type")

	switch dbType {
	case "postgres":
		return setupPostgres()
	default:
		break
	}

	return &sql.DB{}, errors.New("unsupported db: " + dbType)
}

func setupPostgres() (*sql.DB, error) {
	var (
		host   = viper.GetString("db.host")
		port   = viper.GetInt64("db.port")
		user   = viper.GetString("db.user")
		pass   = viper.GetString("db.pass")
		dbname = viper.GetString("db.dbname")
		ssl    = viper.GetString("db.sslmode")
	)

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, dbname, ssl)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		return &sql.DB{}, err
	}

	err = db.Ping()
	if err != nil {
		return &sql.DB{}, err
	}

	log.Info("connected to the postgres database")
	return db, nil
}
