package query

import (
	"database/sql"
	"fmt"
	"log"
	"seraphim/config"

	_ "github.com/go-sql-driver/mysql"
)

func FetchTablesForDb(db string, conn config.StoredConnection) []string {
	tables := make([]string, 0)
	switch conn.Provider {
	case "mysql":
		db, err := sql.Open(conn.Provider, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", conn.User, conn.Password, conn.Host, conn.Port, db))
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Ping the database to check the connection
		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}

		// Query to get the list of databases
		rows, err := db.Query("SHOW TABLES")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var tableName string
			err := rows.Scan(&tableName)
			if err != nil {
				log.Fatal(err)
			}
			tables = append(tables, tableName)
		}

		// Check for errors from iterating over rows
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}

	case "postgres":

	default:
		log.Fatal("Unknown provider")
	}
	return tables
}

func FetchDbList(conn config.StoredConnection) []string {
	dbs := make([]string, 0)
	switch conn.Provider {
	case "mysql":
		db, err := sql.Open(conn.Provider, fmt.Sprintf("%s:%s@tcp(%s:%d)/", conn.User, conn.Password, conn.Host, conn.Port))
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		// Ping the database to check the connection
		err = db.Ping()
		if err != nil {
			log.Fatal(err)
		}

		// Query to get the list of databases
		rows, err := db.Query("SHOW DATABASES")
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		for rows.Next() {
			var dbName string
			err := rows.Scan(&dbName)
			if err != nil {
				log.Fatal(err)
			}
			dbs = append(dbs, dbName)
		}

		// Check for errors from iterating over rows
		if err = rows.Err(); err != nil {
			log.Fatal(err)
		}

	case "postgres":

	default:
		log.Fatal("Unknown provider")
	}
	return dbs
}
