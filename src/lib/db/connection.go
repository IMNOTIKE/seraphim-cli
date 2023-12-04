package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DatabaseConfig struct {
	DatabaseName string
	Host         string
	Port         int
	User         string
	Password     string
}

func ConnectToPostgreSQL(config DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=disable",
		config.User, config.Password, config.DatabaseName, config.Host, config.Port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectToMySQL(config DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.User, config.Password, config.Host, config.Port, config.DatabaseName)
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func ConnectToSQLite(config DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", config.DatabaseName)
	if err != nil {
		return nil, err
	}
	return db, nil
}
