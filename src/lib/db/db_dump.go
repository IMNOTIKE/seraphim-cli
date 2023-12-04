package db

import (
	"fmt"
	"os"
	"seraphim/config"
	"seraphim/lib/bubble/selector"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DbDumpModel struct {
	List         list.Model
	DelegateKeys *delegateKeyMap
	Spinner      spinner.Model
	Err          error

	AvailableConnections      []config.StoredConnection
	SelectedConnectionDetails config.StoredConnection
	Database                  string
	Choosing                  bool
	Tables                    []string
}

type SelectSuccessMsg struct {
	Err       error
	ResultSet []string
}

func (dbm DbDumpModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

var (
	appStyle = lipgloss.NewStyle().Padding(1, 2)

	titleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#25A065")).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render
)

type listItem struct {
	tag  string
	host string
}

func (i listItem) Title() string       { return i.tag }
func (i listItem) Description() string { return i.host }
func (i listItem) FilterValue() string { return i.tag }

func (dbm DbDumpModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		dbm.List.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return dbm, tea.Quit
		case "enter":
			dbm.Choosing = false
			// get selected stored connection

			return dbm, tea.Batch(dbm.FetchTableList(config.StoredConnection{}), spinner.Tick)
		}
	case SelectSuccessMsg:
		dbm.Choosing = true
		dbm.Tables = msg.ResultSet
		return dbm, tea.Quit
	}

	if dbm.Choosing {
		var cmd tea.Cmd
		dbm.List, cmd = dbm.List.Update(msg)
		return dbm, cmd
	}

	return dbm, nil
}

func (dbm DbDumpModel) View() string {
	if dbm.Choosing {
		return fmt.Sprintf("Select a stored connection: \n%s", dbm.List.View())
	}

	if err := dbm.Err; err != nil {
		return fmt.Sprintf("Sorry, could not fetch tables: \n%s", err)
	}

	return "Press Ctrl+C to Exit"
}

func (dbm DbDumpModel) FetchTableList(dbConfig config.StoredConnection) tea.Cmd {
	return func() tea.Msg {
		dbm.Choosing = false
		return SelectSuccessMsg{
			Err:       nil,
			ResultSet: []string{"test", "temp"},
		}
	}
}

func RunDumpCommand(config *config.SeraphimConfig) {
	numItems := len(config.Stored_Connections)
	items := make([]list.Item, numItems)
	delegateKeys := newDelegateKeyMap()
	var i int
	for _, m := range config.Stored_Connections {
		for key, value := range m {
			items[i] = listItem{
				tag:  key,
				host: value.Host,
			}
		}
		i++
	}
	delegate := newItemDelegate(delegateKeys)
	StoredConnectionList := list.New(items, delegate, 0, 0)
	StoredConnectionList.SetShowFilter(true)
	StoredConnectionList.SetShowTitle(false)
	StoredConnectionList.Styles.Title = titleStyle

	s := spinner.New()
	s.Spinner = spinner.Dot

	initialModel := DbDumpModel{
		List:     StoredConnectionList,
		Spinner:  s,
		Choosing: true,
		Tables:   nil,
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	model, err := p.Run()
	if err != nil {
		fmt.Printf("FATAL -- Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	tables := model.(DbDumpModel).Tables
	selectedDb := model.(DbDumpModel).SelectedConnectionDetails
	if tables != nil {
		selector.RunTableSelector(selectedDb, tables)
	} else {
		fmt.Println("Tables were not set")
	}

}
