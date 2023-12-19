package db

import (
	"fmt"
	"seraphim/lib/config"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type RcModel struct {
	StoredConnectionsList list.Model
	DelegateKeys          *delegateKeyMap
	Err                   error

	SelectedConnections []config.StoredConnection
	ChoosingConnections bool
	Done                bool
}

func (rcm RcModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (rcm RcModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		rcm.StoredConnectionsList.SetSize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return rcm, tea.Quit
		case " ":
			// Handle multiselect
		case "enter":
			// pass to huh.Confirm
		}
	}
	var cmd tea.Cmd
	rcm.StoredConnectionsList, cmd = rcm.StoredConnectionsList.Update(msg)
	return rcm, cmd
}

func (rcm RcModel) View() string {

	s := "Press Ctrl+C to Exit"

	if rcm.ChoosingConnections {
		s = fmt.Sprintf("Select a stored connection: \n%s", rcm.StoredConnectionsList.View())
	}

	if err := rcm.Err; err != nil {
		return fmt.Sprintf("Sorry, could not fetch tables: \n%s", err)
	}

	return s
}
