/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

// databaseCmd represents the database command
var databaseCmd = &cobra.Command{
	Use:     "database",
	Aliases: []string{"db"},
	Short:   "Common interface for multiple datasources",
	Long: `This command aims at providing the user a common interface for 
	many different databases`,
}

func init() {
	rootCmd.AddCommand(databaseCmd)
}
