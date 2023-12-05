package db

import (
	"fmt"
	"os"
	"seraphim/lib/config"
	qh "seraphim/lib/db/query"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	btea "github.com/charmbracelet/bubbletea"
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
	DumpPathInput         textinput.Model
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
	return btea.Batch(textinput.Blink, btea.EnterAltScreen)
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
	listDelegate   list.DefaultDelegate
)

//------------------------------------------------//
//
//           List items declarations
//
//------------------------------------------------//

type ConnListItem struct {
	tag  string
	host string
	user string
}

func (c ConnListItem) Title() string       { return c.tag }
func (c ConnListItem) Description() string { return c.user + "@" + c.host }
func (c ConnListItem) FilterValue() string { return c.tag }

type DbListItem struct {
	name string
}

func (d DbListItem) Title() string       { return d.name }
func (d DbListItem) Description() string { return "" }
func (d DbListItem) FilterValue() string { return d.name }

type TableListItem struct {
	name string
}

func (t TableListItem) Title() string       { return t.name }
func (t TableListItem) Description() string { return "" }
func (t TableListItem) FilterValue() string { return t.name }

func (dbm DbDumpModel) updateConnChosingView(msg btea.Msg) (btea.Model, btea.Cmd) {

	switch msg := msg.(type) {
	case btea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		dbm.StoredConnectionsList.SetSize(msg.Width-h, msg.Height-v)
	case btea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return dbm, btea.Quit
		case "enter":
			dbm.ChoosingConnection = false
			selectedItem := dbm.StoredConnectionsList.SelectedItem().(ConnListItem)
			for _, conn := range seraphimConfig.Stored_Connections {
				t := conn[selectedItem.tag]
				dbm.SelectedConnectionDetails = t
				break
			}
			dbs := qh.FetchDbList(dbm.SelectedConnectionDetails)
			dbsListItems := make([]list.Item, len(dbs))
			for i, db := range dbs {
				dbsListItems[i] = DbListItem{
					name: db,
				}
			}
			DatabasesList := list.New(dbsListItems, listDelegate, 0, 0)
			DatabasesList.SetShowFilter(true)
			DatabasesList.SetShowTitle(false)
			DatabasesList.Styles.Title = titleStyle
			dbm.DatabasesList = DatabasesList
			dbm.ChoosingDatabases = true
			return dbm, btea.ClearScreen
		}
	}
	var cmd btea.Cmd
	dbm.StoredConnectionsList, cmd = dbm.StoredConnectionsList.Update(msg)
	return dbm, cmd
}

func (dbm DbDumpModel) updateDbChosingView(msg btea.Msg) (btea.Model, btea.Cmd) {
	switch msg := msg.(type) {
	case btea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		dbm.DatabasesList.SetSize(msg.Width-h, msg.Height-v)
	case btea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return dbm, btea.Quit
		case "backspace":
			dbm.ChoosingConnection = true
			dbm.ChoosingDatabases = false
		case "enter":
			dbm.ChoosingDatabases = false
			tables := qh.FetchTablesForDb("", dbm.SelectedConnectionDetails)
			tableListItems := make([]list.Item, len(tables))
			for i, table := range tables {
				tableListItems[i] = TableListItem{
					name: table,
				}
			}
			TablesList := list.New(tableListItems, listDelegate, 0, 0)
			TablesList.SetShowFilter(true)
			TablesList.SetShowTitle(false)
			TablesList.Styles.Title = titleStyle
			dbm.TablesList = TablesList
			dbm.ChoosingTables = true
			return dbm, btea.ClearScreen
		}
	}
	var cmd btea.Cmd
	dbm.DatabasesList, cmd = dbm.DatabasesList.Update(msg)
	return dbm, cmd
}

func (dbm DbDumpModel) updateTableChosingView(msg btea.Msg) (btea.Model, btea.Cmd) {
	switch msg := msg.(type) {
	case btea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		dbm.TablesList.SetSize(msg.Width-h, msg.Height-v)
	case btea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return dbm, btea.Quit
		case "backspace":
			dbm.ChoosingDatabases = true
			dbm.ChoosingTables = false
		case "enter":
			dbm.ChoosingTables = false
			dbm.TypingPath = true
			return dbm, btea.ClearScreen
		}
	}
	var cmd btea.Cmd
	dbm.TablesList, cmd = dbm.TablesList.Update(msg)
	return dbm, cmd
}

func (dbm DbDumpModel) Update(msg btea.Msg) (btea.Model, btea.Cmd) {

	if dbm.ChoosingConnection {
		return dbm.updateConnChosingView(msg)
	}
	if dbm.ChoosingDatabases {
		return dbm.updateDbChosingView(msg)
	}
	if dbm.ChoosingTables {
		return dbm.updateTableChosingView(msg)
	}

	if dbm.TypingPath {
		var cmd btea.Cmd
		dbm.DumpPathInput, cmd = dbm.DumpPathInput.Update(msg)
		dbm.DumpPathInput.Focus()
		return dbm, cmd
	}

	return dbm, nil
}

func (dbm DbDumpModel) View() string {

	if dbm.ChoosingConnection {
		return fmt.Sprintf("Select a stored connection: \n%s", dbm.StoredConnectionsList.View())
	}

	if dbm.ChoosingDatabases {
		return fmt.Sprintf("Select a one or more databases: \n%s", dbm.DatabasesList.View())
	}

	if dbm.ChoosingTables {
		return fmt.Sprintf("Select a one or more tables: \n%s", dbm.TablesList.View())
	}

	if dbm.TypingPath {
		return fmt.Sprintf("Select a one or more tables: \n%s", dbm.DumpPathInput.View())
	}

	if err := dbm.Err; err != nil {
		return fmt.Sprintf("Sorry, could not fetch tables: \n%s", err)
	}

	return "Press Ctrl+C to Exit"
}

func RunDumpCommand(config *config.SeraphimConfig) {
	seraphimConfig = *config
	numItems := len(config.Stored_Connections)
	delegateKeys := newDelegateKeyMap()
	items := make([]list.Item, numItems)
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
	listDelegate = delegate

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
