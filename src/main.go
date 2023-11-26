package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Seraphim CLI tool"
	app.Usage = "A simple CLI application containing all sorts of useful tools"

	dbDumpFlags := []cli.Flag{
		cli.BoolFlag{
			Name: "host",
			Usage: "Set the host of the database",
			Required: true,
		},
	}

	dbSubCommands := []cli.Command{
		{
			Name: "dump",
			Aliases: []string{"d"},
			Usage: "Create a dump from a database",
			Category: "database",
			Flags: dbDumpFlags,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "database",
			Aliases: []string{"db"},
			Usage:   "Database tool-belt",
			Category: "database",
			Subcommands: dbSubCommands,
			Action:  runBubbleTea,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func runBubbleTea(c *cli.Context) error {
	p := progress.New()

	program := tea.NewProgram(p)
	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}