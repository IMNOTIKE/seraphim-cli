/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"seraphim/lib/config"

	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the configuration from the cli",
	Run: func(cmd *cobra.Command, args []string) {
		r := config.RunCfgEditForm(&seraphimConfig)
		if r.Err != nil {
			log.Fatal("something went wrong")
		}
		if operationResult := config.SaveConfig(r.EditedConfig); operationResult.Err == nil {
			fmt.Printf("%s\n", operationResult.Msg)
		} else {
			fmt.Printf("Oh no, something went wrong: \n%v", operationResult.Err.Error())
		}
	},
}

func init() {
	configCmd.AddCommand(editCmd)
}
