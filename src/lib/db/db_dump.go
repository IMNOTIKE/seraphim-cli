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
			var dbsListItems []list.Item
			for _, db := range dbs {
				dbsListItems = append(dbsListItems, DbListItem{
					name: db,
				})
			}
			dbm.DatabasesList.SetItems(dbsListItems)
			dbm.ChoosingDatabases = true
			return dbm, nil
		}
	}
	return dbm, nil
}

func (dbm DbDumpModel) updateDbChosingView(msg btea.Msg) (btea.Model, btea.Cmd) {
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
			tables := qh.FetchTablesForDb("", dbm.SelectedConnectionDetails)
			var tableListItems []list.Item
			for _, table := range tables {
				tableListItems = append(tableListItems, DbListItem{
					name: table,
				})
			}
			dbm.TablesList.SetItems(tableListItems)
			dbm.ChoosingDatabases = true
			return dbm, nil
		}
	}
	return dbm, nil
}

func (dbm DbDumpModel) updateTableChosingView(msg btea.Msg) (btea.Model, btea.Cmd) {
	switch msg := msg.(type) {
	case btea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		dbm.StoredConnectionsList.SetSize(msg.Width-h, msg.Height-v)
	case btea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return dbm, btea.Quit
		case "enter":
			dbm.ChoosingTables = false
			dbm.TypingPath = true
			return dbm, nil
		}
	}
	return dbm, nil
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

	if dbm.ChoosingConnection {
		var cmd btea.Cmd
		dbm.StoredConnectionsList, cmd = dbm.StoredConnectionsList.Update(msg)
		return dbm, cmd
	}

	if dbm.ChoosingDatabases {
		var cmd btea.Cmd
		dbm.DatabasesList, cmd = dbm.DatabasesList.Update(msg)
		return dbm, cmd
	}

	if dbm.ChoosingTables {
		var cmd btea.Cmd
		dbm.TablesList, cmd = dbm.TablesList.Update(msg)
		return dbm, cmd
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

	DatabasesList := list.New([]list.Item{}, delegate, 0, 0)
	TablesList := list.New([]list.Item{}, delegate, 0, 0)

	s := spinner.New()
	s.Spinner = spinner.Dot

	initialModel := DbDumpModel{
		StoredConnectionsList: StoredConnectionList,
		DatabasesList:         DatabasesList,
		TablesList:            TablesList,
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
