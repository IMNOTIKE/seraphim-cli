/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"seraphim/lib/db"

	"github.com/spf13/cobra"
)

// editConnectionCmd represents the editConnection command
var editConnectionCmd = &cobra.Command{
	Use:     "edit-connection",
	Short:   "Allows the user to select a stored connection and edit it from the command line",
	Aliases: []string{"ec"},
	Run: func(cmd *cobra.Command, args []string) {
		db.RunStoredConnectionEditHandler(&seraphimConfig)
	},
}

func init() {
	databaseCmd.AddCommand(editConnectionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// editConnectionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// editConnectionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
