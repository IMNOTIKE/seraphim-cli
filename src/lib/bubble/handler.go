package bubble

import (
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/urfave/cli"
)

func runBubbleTea(c *cli.Context) error {
	p := progress.New()

	program := tea.NewProgram(p)
	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}
