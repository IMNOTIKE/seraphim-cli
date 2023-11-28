package db

import (
	"fmt"
	"os"
	"seraphim/config"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

type DbDumpModel struct {
	Input   textinput.Model
	Spinner spinner.Model
	Err     error

	// Add states for various operations
	AvailableConnections      []config.StoredConnection
	SelectedConnectionDetails config.StoredConnection
	Database                  string
	Typing                    bool
	Loading                   bool
	Tables                    []string
}

type SelectSuccessMsg struct {
	Err       error
	ResultSet []string
}

func (dbm DbDumpModel) Init() tea.Cmd {
	return textinput.Blink
}

func (dbm DbDumpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return dbm, tea.Quit
		case "enter":
			dbm.Typing = false
			dbm.Loading = true
			return dbm, tea.Batch(dbm.FetchTableList(), spinner.Tick)
		}
	case SelectSuccessMsg:
		dbm.Loading = false
		dbm.Typing = true
		return dbm, nil
	}

	if dbm.Loading {
		var cmd tea.Cmd
		dbm.Spinner, cmd = dbm.Spinner.Update(msg)
		return dbm, cmd
	}

	if dbm.Typing {
		var cmd tea.Cmd
		dbm.Input, cmd = dbm.Input.Update(msg)
		return dbm, cmd
	}

	return dbm, nil
}

func (dbm DbDumpModel) View() string {
	if dbm.Typing {
		return fmt.Sprintf("Enter databse name: \n%s", dbm.Input.View())
	}

	if dbm.Loading {
		return fmt.Sprintf("%s Creating dump", dbm.Spinner.View())
	}

	if err := dbm.Err; err != nil {
		return fmt.Sprintf("Sorry, could not fetch tables: \n%s", err)
	}

	return "Press Ctrl+C to Exit"
}

func (dbm DbDumpModel) FetchTableList() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		return SelectSuccessMsg{
			Err:       nil,
			ResultSet: []string{},
		}
	}
}

func RunDumpCommand(config *config.SeraphimConfig) {

	i := textinput.New()
	i.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot

	initialModel := DbDumpModel{
		Input:   i,
		Spinner: s,
		Typing:  true,
	}

	p := tea.NewProgram(initialModel)
	if _, err := p.Run(); err != nil {
		log.Info("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
