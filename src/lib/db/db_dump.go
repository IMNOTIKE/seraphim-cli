package db

import (
	"fmt"
	"log"
	"os"
	"seraphim/lib/bubble/selector"
	"seraphim/lib/config"
	"seraphim/lib/db/query"
	qh "seraphim/lib/db/query"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	btea "github.com/charmbracelet/bubbletea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

/**
* Flow
* --> dumpCmd --> RunDumpCommand --> RunDbSelector --> RunMultiSelectList(Tables) --> CreateDump
* List Selector for connections
* List Selector for DB
* List MultiSelector for Table
* Input for path
* Banner to show result
 */

type DbDumpModel struct {
	StoredConnectionsList list.Model
	DatabasesList         list.Model
	TablesList            list.Model
	DelegateKeys          *delegateKeyMap
	Spinner               spinner.Model
	DbInput               textinput.Model
	Err                   error

	AvailableConnections      []config.StoredConnection
	AvailableDatabase         []string
	AvailableTables           []string
	SelectedConnectionDetails config.StoredConnection
	SelectedDatabases         []string
	SelectedTables            []string
	ChoosingConnection        bool
	ChoosingDatabases         bool
	ChoosingTables            bool
	TypingPath                bool
	Done                      bool
}

func (dbm DbDumpModel) Init() btea.Cmd {
	return btea.EnterAltScreen
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

	seraphimConfig config.SeraphimConfig
)

type ConnListItem struct {
	tag  string
	host string
	user string
}

func (i ConnListItem) Title() string       { return i.tag }
func (i ConnListItem) Description() string { return i.user + "@" + i.host }
func (i ConnListItem) FilterValue() string { return i.tag }

type DbListItem struct {
	name string
}

func (i DbListItem) Title() string       { return i.name }
func (i DbListItem) Description() string { return "" }
func (i DbListItem) FilterValue() string { return i.name }

type TableListItem struct {
	name string
}

func (i TableListItem) Title() string       { return i.name }
func (i TableListItem) Description() string { return "" }
func (i TableListItem) FilterValue() string { return i.name }

func (dbm DbDumpModel) updateConnChosingView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case btea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		dbm.StoredConnectionsList.SetSize(msg.Width-h, msg.Height-v)
	case btea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return dbm, btea.Quit
		case "enter":
			dbm.ChoosingConnection = false
			selectedItem := dbm.StoredConnectionsList.SelectedItem().(ConnListItem)
			for _, conn := range seraphimConfig.Stored_Connections {
				t := conn[selectedItem.tag]
				dbm.SelectedConnectionDetails = t
				break
			}
			dbm.ChoosingDatabases = true
			return dbm, nil
		}
	}
	return dbm, tea.Quit
}

func (dbm DbDumpModel) updateDbChosingView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case btea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		dbm.StoredConnectionsList.SetSize(msg.Width-h, msg.Height-v)
	case btea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return dbm, btea.Quit
		case "enter":
			dbm.ChoosingConnection = false
			selectedItem := dbm.StoredConnectionsList.SelectedItem().(ConnListItem)
			for _, conn := range seraphimConfig.Stored_Connections {
				t := conn[selectedItem.tag]
				dbm.SelectedConnectionDetails = t
				break
			}
			dbm.AvailableDatabase = query.FetchDbList(dbm.SelectedConnectionDetails)
			dbm.ChoosingDatabases = true
			return dbm, nil
		}
	}
	return dbm, tea.Quit
}

func (dbm DbDumpModel) updateTableChosingView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case btea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		dbm.StoredConnectionsList.SetSize(msg.Width-h, msg.Height-v)
	case btea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return dbm, btea.Quit
		case "enter":
			dbm.ChoosingDatabases = false
			// Get selected dbs
			dbm.AvailableTables = query.FetchTablesForDb("", dbm.SelectedConnectionDetails)
			dbm.ChoosingTables = true
			return dbm, nil
		}
	}
	return dbm, tea.Quit
}

func (dbm DbDumpModel) Update(msg btea.Msg) (btea.Model, btea.Cmd) {
	if dbm.ChoosingConnection {
		dbm.updateConnChosingView(msg)
	}
	if dbm.ChoosingDatabases {
		dbm.updateDbChosingView(msg)
	}
	if dbm.ChoosingTables {
		dbm.updateTableChosingView(msg)
	}
	return dbm, nil
}

func (dbm DbDumpModel) View() string {
	if dbm.ChoosingConnection {
		return fmt.Sprintf("Select a stored connection: \n%s", dbm.StoredConnectionsList.View())
	}

	if err := dbm.Err; err != nil {
		return fmt.Sprintf("Sorry, could not fetch tables: \n%s", err)
	}

	return "Press Ctrl+C to Exit"
}

func (dbm DbDumpModel) PerformDump(dbConfig config.StoredConnection, defDmpPath string) {
	dbs := qh.FetchDbList(dbConfig)
	selectDbResult := selector.RunDbSelector(dbConfig, dbs)
	if selectDbResult.Err != nil {
		log.Fatal(selectDbResult.Err.Error())
	}
	db := selectDbResult.Result[0]
	tables := qh.FetchTablesForDb(db, dbConfig)
	selector.RunMultiSelectList(tables, dbConfig, db, defDmpPath)
}

func RunDumpCommand(config *config.SeraphimConfig) {
	seraphimConfig = *config
	numItems := len(config.Stored_Connections)
	items := make([]list.Item, numItems)
	delegateKeys := newDelegateKeyMap()
	var i int
	for _, m := range config.Stored_Connections {
		for key, value := range m {
			items[i] = ConnListItem{
				tag:  key,
				host: value.Host,
				user: value.User,
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
		StoredConnectionsList: StoredConnectionList,
		Spinner:               s,
		ChoosingConnection:    true,
	}

	p := btea.NewProgram(initialModel, btea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		fmt.Printf("FATAL -- Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
