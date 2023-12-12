/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"seraphim/lib/config"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Initialize configuration file in the default directory ($HOME/.config/seraphim/)",
	Run: func(cmd *cobra.Command, args []string) {
		if res := config.InitConfig(); res.Err == nil {
			fmt.Println(res.Msg)
		} else {
			fmt.Println(res.Err)
		}
	},
}

func init() {
	configCmd.AddCommand(initCmd)
}
