package db

import (
	"fmt"
	"os"
	"seraphim/config"
	"time"

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
	Loading                   bool
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
			dbm.Loading = true
			// get selected stored connection
			return dbm, tea.Batch(dbm.FetchTableList(config.StoredConnection{}), spinner.Tick)
		}
	case SelectSuccessMsg:
		dbm.Loading = false
		dbm.Choosing = true
		return dbm, nil
	}

	if dbm.Loading {
		var cmd tea.Cmd
		dbm.Spinner, cmd = dbm.Spinner.Update(msg)
		return dbm, cmd
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

	if dbm.Loading {
		return fmt.Sprintf("%s Creating dump", dbm.Spinner.View())
	}

	if err := dbm.Err; err != nil {
		return fmt.Sprintf("Sorry, could not fetch tables: \n%s", err)
	}

	return "Press Ctrl+C to Exit"
}

func (dbm DbDumpModel) FetchTableList(dbConfig config.StoredConnection) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		dbm.Choosing = false
		dbm.Loading = false
		return SelectSuccessMsg{
			Err:       nil,
			ResultSet: []string{},
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
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("FATAL -- Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
