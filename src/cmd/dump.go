/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"seraphim/config"
	"seraphim/lib/db"

	"database/sql"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/JamesStewy/go-mysqldump"
	_ "github.com/go-sql-driver/mysql"
	pg_dumper "github.com/habx/pg-commands"
)

// dumpCmd represents the dump command
var dumpCmd = &cobra.Command{
	Use:     "dump",
	Aliases: []string{"dmp"},
	Short:   "Create a database dump",
	Long:    `Create a dump of the selected database`,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		db.RunDumpCommand(&seraphimConfig)
	},
}

func init() {
	databaseCmd.AddCommand(dumpCmd)
}

func CreateDump(selected config.StoredConnection, dumpPath string, selectedDb string) {
	username := selected.User
	password := selected.Password
	hostname := selected.Host
	port := selected.Port
	dbname := selectedDb
	driver := selected.Provider

	dumpDir := "dumps"                                              // you should create this directory
	dumpFilenameFormat := fmt.Sprintf("%s-20060102T150405", dbname) // accepts time layout string and add .sql at the end of file

	switch driver {
	case "mysql":
		db, err := sql.Open(driver, fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", username, password, hostname, port, dbname))
		if err != nil {
			fmt.Println("Error opening database: ", err)
			return
		}

		// Register database with mysqldump
		dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
		if err != nil {
			fmt.Println("Error registering databse:", err)
			return
		}

		// Dump database to file
		resultFilename, err := dumper.Dump()
		if err != nil {
			fmt.Println("Error dumping:", err)
			return
		}
		fmt.Printf("File is saved to %s", resultFilename)

		// Close dumper and connected database
		dumper.Close()
	case "postgres":
		dump, _ := pg_dumper.NewDump(&pg_dumper.Postgres{
			Host:     hostname,
			Port:     port,
			DB:       dbname,
			Username: username,
			Password: password,
		})
		dumpExec := dump.Exec(pg_dumper.ExecOptions{StreamPrint: false})
		if dumpExec.Error != nil {
			fmt.Println(dumpExec.Error.Err)
			fmt.Println(dumpExec.Output)

		} else {
			fmt.Println("Dump success")
			fmt.Println(dumpExec.Output)
		}
	default:
		log.Fatal("Unknown driver")
	}
}

// func RestoreDB(){
// 	restore, _ := pg_dumper.NewRestore(&pg_dumper.Postgres{
// 		Host:     "localhost",
// 		Port:     5432,
// 		DB:       "dev_example",
// 		Username: "example",
// 		Password: "example",
// 	})
// 	restoreExec := restore.Exec(dumpExec.File, pg_dumper.ExecOptions{StreamPrint: false})
// 	if restoreExec.Error != nil {
// 		fmt.Println(restoreExec.Error.Err)
// 		fmt.Println(restoreExec.Output)

// 	} else {
// 		fmt.Println("Restore success")
// 		fmt.Println(restoreExec.Output)

// 	}
// }
