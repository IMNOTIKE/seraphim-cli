/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"seraphim/lib/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var seraphimConfig config.SeraphimConfig
var versionRequested bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "seraphim",
	Short: "Modular and varied toolbelt",
	Long:  "Serpahim aims at providing the user with several commands\nto make life easieer\nRequired dependencies:\n- mysqldump",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Thank you for using seraphim")
		if versionRequested {
			fmt.Printf("\u251C\u279D  App: %s\n", seraphimConfig.Branding.Name)
			fmt.Printf("\u2514\u279D  Version: %s\n", seraphimConfig.Version)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/seraphim/seraphim.yaml)")
	rootCmd.Flags().BoolVarP(&versionRequested, "version", "v", false, "Application version")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".seraphim" (without extension).
		viper.AddConfigPath(home + "/.config/seraphim")
		viper.SetConfigType("yaml")
		viper.SetConfigName("seraphim")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		viper.Unmarshal(&seraphimConfig)
	}
}

func DeleteKeyFromConfig(key string) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".seraphim" (without extension).
		viper.AddConfigPath(home + "/.config/seraphim")
		viper.SetConfigType("yaml")
		viper.SetConfigName("seraphim")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("error")
	}

}
