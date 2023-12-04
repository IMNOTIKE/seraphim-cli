/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"seraphim/lib/db"

	"github.com/spf13/cobra"

	_ "github.com/go-sql-driver/mysql"
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
