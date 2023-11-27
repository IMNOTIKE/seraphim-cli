package main

import "seraphim-cli/lib/store"

// func main() {
// 	app := cli.NewApp()
// 	app.Name = "Seraphim CLI tool"
// 	app.Usage = "A simple CLI application containing all sorts of useful tools"

// 	dbDumpFlags := []cli.Flag{
// 		cli.BoolFlag{
// 			Name: "host",
// 			Usage: "Set the host of the database",
// 			Required: true,
// 		},
// 	}

// 	dbSubCommands := []cli.Command{
// 		{
// 			Name: "dump",
// 			Aliases: []string{"d"},
// 			Usage: "Create a dump from a database",
// 			Category: "database",
// 			Flags: dbDumpFlags,
// 		},
// 	}

// 	app.Commands = []cli.Command{
// 		{
// 			Name:    "database",
// 			Aliases: []string{"db"},
// 			Usage:   "Database tool-belt",
// 			Category: "database",
// 			Subcommands: dbSubCommands,
// 			Action:  runBubbleTea,
// 		},
// 	}

// 	err := app.Run(os.Args)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

func main() {
	// postgresConfig := db.DatabaseConfig{
	// 	Name:     "my_postgres_db",
	// 	Host:     "localhost",
	// 	Port:     5432,
	// 	User:     "postgres",
	// 	Password: "password",
	// }

	// mysqlConfig := db.DatabaseConfig{
	// 	Name:     "my_mysql_db",
	// 	Host:     "localhost",
	// 	Port:     3306,
	// 	User:     "root",
	// 	Password: "password",
	// }

	// sqliteConfig := db.DatabaseConfig{
	// 	Name: "my_sqlite_db.db",
	// }

	// postgresDB, err := db.ConnectToPostgreSQL(postgresConfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// mysqlDB, err := db.ConnectToMySQL(mysqlConfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// sqliteDB, err := db.ConnectToSQLite(sqliteConfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Now you have connections to all your databases: postgresDB, mysqlDB, and sqliteDB
	store.DoesStoreDBFileExist()
}
