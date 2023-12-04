/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"seraphim/lib/db"

	"github.com/spf13/cobra"
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
