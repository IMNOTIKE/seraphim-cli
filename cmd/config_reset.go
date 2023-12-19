/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"seraphim/lib/config"

	"github.com/spf13/cobra"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the configuration file to default settings",
	Run: func(cmd *cobra.Command, args []string) {
		if res := config.ResetConfig(); res.Err == nil {
			fmt.Println(res.Msg)
		} else {
			fmt.Println(res.Err)
		}
	},
}

func init() {
	configCmd.AddCommand(resetCmd)
}
