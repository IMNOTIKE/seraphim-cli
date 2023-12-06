package query

import (
	"database/sql"
	"fmt"
	"log"
	"os/exec"
	"seraphim/lib/config"
	"seraphim/lib/util"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	pgcommands "github.com/habx/pg-commands"
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

func CreateDump(selected config.StoredConnection, dumpPath string, selectedDbs []util.DbListItem, selectedTables []util.TableListItem) bool {
	username := selected.User
	password := selected.Password
	hostname := selected.Host
	port := selected.Port
	driver := selected.Provider

	dumpDir := dumpPath
	dumpFilenameFormat := fmt.Sprintf("%s-%v.sql", "dump", time.Now().Unix())
	switch driver {
	case "mysql":
		// I'd rather use a golang library to avoid external dependencies but
		// no library offers the same flexibility
		sql := fmt.Sprintf("mysqldump -u %s -p%s ", username, password)
		isFirst := true
		if len(selectedDbs) > 0 {
			for _, db := range selectedDbs {
				var b strings.Builder
				b.WriteString(sql)
				b.WriteString(" " + db.Name)
				if len(selectedTables) > 0 {
					for _, v := range selectedTables {
						if v.Db == db.Name {
							if v.Name == "All" {
								b.Reset()
								b.WriteString(sql)
								b.WriteString(" " + db.Name)
								break
							}
							b.WriteString(" " + v.Name)
						}
					}
				}
				if isFirst {
					b.WriteString(fmt.Sprintf(" > %s", dumpFilenameFormat))
					isFirst = false
				} else {
					b.WriteString(fmt.Sprintf(" >> %s", dumpFilenameFormat))
				}
				cmd := exec.Command("bash", "-c", b.String())
				cmd.Dir = dumpDir
				err := cmd.Run()
				if err != nil {
					serr := err.Error()
					fmt.Println(serr)
					return false
				}
			}
		}
		return true
	case "postgres":
		dump, _ := pgcommands.NewDump(&pgcommands.Postgres{
			Host:     hostname,
			Port:     port,
			DB:       "",
			Username: username,
			Password: password,
		})
		dumpExec := dump.Exec(pgcommands.ExecOptions{StreamPrint: false})
		if dumpExec.Error != nil {
			fmt.Println(dumpExec.Error.Err)
			fmt.Println(dumpExec.Output)
			return false

		} else {
			fmt.Println("Dump success")
			fmt.Println(dumpExec.Output)
			return true
		}
	default:
		log.Fatal("Unknown driver")
		return false
	}
}
