/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"seraphim/lib/config"
	"seraphim/lib/db"

	"github.com/spf13/cobra"
)

// addConnectionCmd represents the addConnection command
var addConnectionCmd = &cobra.Command{
	Use:     "add-connection",
	Aliases: []string{"ac"},
	Short:   "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		r := db.RunAdcForm()
		if r.Err != nil {
			log.Fatal("something went wrong")
		}
		if operationResult := config.AddConnection(false, config.SeraphimConfig{}, r.NewConnection, r.Tag); operationResult.Err == nil {
			fmt.Printf("%s\n", operationResult.Msg)
		} else {
			fmt.Printf("Oh no, something went wrong: \n%v", operationResult.Err.Error())
		}
	},
}

func init() {
	databaseCmd.AddCommand(addConnectionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addConnectionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addConnectionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
