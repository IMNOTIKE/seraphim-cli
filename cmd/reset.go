/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the configuration file to default settings",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("reset called")
	},
}

func init() {
	configCmd.AddCommand(resetCmd)
}
