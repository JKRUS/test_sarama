package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

//var DB *sqlx.DB = sqlx.NewDb(db, "postgres")

func NewPostgresDB() error {
	cfg, err := newConfig()
	if err != nil {
		return err
	}
	DB, err = sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.UserName, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}
	return nil
}
