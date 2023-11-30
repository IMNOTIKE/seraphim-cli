package dialog

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type DialogUseCase string

const (
	DeleteStoredConnection DialogUseCase = "deleteStoredLocation"
)

type confirmationDialog struct {
	input      textinput.Model
	options    []string
	title      string
	interacted bool
	confirmed  bool
}

func (cd confirmationDialog) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

func (cd confirmationDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return cd, nil
}

func (cd confirmationDialog) View() string {
	return ""
}

func RunDialog(dialogMessage string, useCase DialogUseCase, params ...any) {

	switch useCase {
	case DeleteStoredConnection:
		//TODO

	default:
		// TODO
	}

	i := textinput.New()
	i.Focus()
	initialModel := confirmationDialog{
		input:      i,
		options:    make([]string, 0),
		title:      "",
		interacted: false,
		confirmed:  false,
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("FATAL -- Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
