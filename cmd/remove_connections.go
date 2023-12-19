/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// removeConnectionCmd represents the removeConnection command
var removeConnectionCmd = &cobra.Command{
	Use:   "remove-connections",
	Short: "Remove one or more stored connections from config",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("removeConnection called")
	},
}

func init() {
	databaseCmd.AddCommand(removeConnectionCmd)
}
