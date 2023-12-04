package dialog

import (
	"errors"
	"fmt"
	"seraphim/config"
	"slices"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type DialogUseCase string

const (
	DeleteStoredConnection DialogUseCase = "deleteStoredLocation"
	RequestUserInput       DialogUseCase = "requestUserInput"
)

type confirmationDialog struct {
	input   textinput.Model
	err     error
	options string
}

func (cd confirmationDialog) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

type (
	errMsg error
)

var (
	keyIndex     int
	removeParams []string
	question     string
	notDeleted   error
)

func (cd confirmationDialog) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			notDeleted = errors.New("not deleted")
			return cd, tea.Quit
		case tea.KeyEnter:
			if slices.Contains([]string{"y", "yes"}, cd.input.Value()) {
				config.RemoveStoredConnection(keyIndex, removeParams...)
				notDeleted = nil
				return cd, tea.Quit
			}
			notDeleted = errors.New("not deleted")
			return cd, tea.Quit
		}
	// We handle errors just like any other message
	case errMsg:
		cd.err = msg
		notDeleted = msg
		return cd, nil
	}

	cd.input, cmd = cd.input.Update(msg)
	return cd, cmd
}

func (cd confirmationDialog) View() string {
	return fmt.Sprintf(
		question+"\n\n%s\n\n%s",
		cd.input.View(),
		"(ctrl+C to abort)",
	) + "\n"
}

func RunDialog(dialogMessage string, useCase DialogUseCase, index int, params ...string) (string, error) {

	var initialModel confirmationDialog
	keyIndex = index
	removeParams = params
	question = dialogMessage
	i := textinput.New()
	i.Focus()
	switch useCase {
	case DeleteStoredConnection:
		i.CharLimit = 3
		var valid bool
		if len(params) < 1 {
			return "", errors.New("expected params to be one or more, got " + strconv.Itoa(len(params)))
		} else {
			valid = true
		}

		if valid {
			initialModel = confirmationDialog{
				input:   i,
				options: "y/yes or N/No",
			}
			i.Placeholder = initialModel.options
			p := tea.NewProgram(initialModel, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Printf("FATAL -- Alas, there's been an error: %v", err)
				return "", err
			}
			return "", notDeleted
		}
	case RequestUserInput:
		i.CharLimit = 20
		var valid bool
		if len(params) < 1 {
			return "", errors.New("expected params to be one or more, got " + strconv.Itoa(len(params)))
		} else {
			valid = true
		}

		if valid {
			initialModel = confirmationDialog{
				input:   i,
				options: "Insert the path for the dump or press ENTER to use $HOME",
			}
			i.Placeholder = initialModel.options
			p := tea.NewProgram(initialModel, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Printf("FATAL -- Alas, there's been an error: %v", err)
				return "", err
			}
			return i.Value(), notDeleted
		}
	default:
		err := errors.New("unknown dialog variant")
		fmt.Printf("FATAL -- Alas, there's been an error: %v", err)
		return "", err
	}

	return "", errors.New("something went wrong while handling the operation")
}
