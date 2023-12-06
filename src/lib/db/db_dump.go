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
	InputDumpPathValue        string
	ChoosingConnection        bool
	ChoosingDatabases         bool
	ChoosingTables            bool
	TypingPath                bool
	Done                      bool
	End                       bool
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

	pathInputTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("230"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle         = focusedStyle.Copy()
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
	Name     string
	Selected bool
}

func (d DbListItem) Title() string       { return d.Name }
func (d DbListItem) Description() string { return "" }
func (d DbListItem) FilterValue() string { return d.Name }

type TableListItem struct {
	Name     string
	Db       string
	Selected bool
}

func (t TableListItem) Title() string       { return t.Name }
func (t TableListItem) Description() string { return t.Db }
func (t TableListItem) FilterValue() string { return t.Name }

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
				if t != (config.StoredConnection{}) {
					dbm.SelectedConnectionDetails = t
					break
				}
			}
			dbs := qh.FetchDbList(dbm.SelectedConnectionDetails)
			dbsListItems := make([]list.Item, len(dbs))
			for i, db := range dbs {
				dbsListItems[i] = DbListItem{
					Name: db,
				}
			}
			DatabasesList := list.New(dbsListItems, listDelegate, 0, 0)
			DatabasesList.SetShowFilter(true)
			DatabasesList.SetShowTitle(false)
			DatabasesList.Styles.Title = titleStyle
			dbm.DatabasesList = DatabasesList
			dbm.ChoosingDatabases = true
			return dbm, func() btea.Msg {
				return btea.WindowSizeMsg{
					Height: dbm.StoredConnectionsList.Height(),
					Width:  dbm.StoredConnectionsList.Width(),
				}
			}
		}
	}
	var cmd btea.Cmd
	dbm.StoredConnectionsList, cmd = dbm.StoredConnectionsList.Update(msg)
	return dbm, cmd
}

func (dbm DbDumpModel) updateDbChosingView(msg btea.Msg) (btea.Model, btea.Cmd) {
	switch msg := msg.(type) {
	case btea.WindowSizeMsg:
		dbm.DatabasesList.SetSize(msg.Width, msg.Height)
	case btea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return dbm, btea.Quit
		case "backspace":
			dbm.ChoosingConnection = true
			dbm.ChoosingDatabases = false
		case " ":
			// Handle multiple selection
		case "enter":
			dbm.ChoosingDatabases = false
			selectedItem := dbm.DatabasesList.SelectedItem().(DbListItem)
			dbm.SelectedDatabases = []string{selectedItem.Name}
			// Should fetch tables for all selected dbs
			tables := qh.FetchTablesForDb(selectedItem.Name, dbm.SelectedConnectionDetails)
			tableListItems := make([]list.Item, len(tables))
			for i, table := range tables {
				tableListItems[i] = TableListItem{
					Name: table,
				}
			}
			TablesList := list.New(tableListItems, listDelegate, 0, 0)
			TablesList.SetShowFilter(true)
			TablesList.SetShowTitle(false)
			TablesList.Styles.Title = titleStyle
			dbm.TablesList = TablesList
			dbm.ChoosingTables = true
			return dbm, func() btea.Msg {
				return btea.WindowSizeMsg{
					Height: dbm.DatabasesList.Height(),
					Width:  dbm.DatabasesList.Width(),
				}
			}
		}
	}
	var cmd btea.Cmd
	dbm.DatabasesList, cmd = dbm.DatabasesList.Update(msg)
	return dbm, cmd
}

func (dbm DbDumpModel) updateTableChosingView(msg btea.Msg) (btea.Model, btea.Cmd) {
	switch msg := msg.(type) {
	case btea.WindowSizeMsg:
		dbm.TablesList.SetSize(msg.Width, msg.Height)
	case btea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return dbm, btea.Quit
		case " ":
			// Handle multiple selection
		case "backspace":
			dbm.ChoosingDatabases = true
			dbm.ChoosingTables = false
		case "enter":
			dbm.ChoosingTables = false
			selectedTable := dbm.TablesList.SelectedItem().(TableListItem)
			dbm.SelectedTables = []string{selectedTable.Name}
			dbm.TypingPath = true
			return dbm, nil
		}
	}
	var cmd btea.Cmd
	dbm.TablesList, cmd = dbm.TablesList.Update(msg)
	return dbm, cmd
}

func (dbm DbDumpModel) updatePathInputView(msg btea.Msg) (btea.Model, btea.Cmd) {
	switch msg := msg.(type) {
	case btea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		dbm.TablesList.SetSize(msg.Width-h, msg.Height-v)
	case btea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return dbm, btea.Quit
		case "alt+backspace":
			dbm.ChoosingTables = true
			dbm.TypingPath = false
		case "enter":
			dbm.TypingPath = false
			inputPath := dbm.DumpPathInput.Value()
			if inputPath == "" {
				inputPath = seraphimConfig.Default_dump_path
			}
			dbm.InputDumpPathValue = inputPath
			dbm.Done = true
			return dbm, nil
		}
	}
	var cmd btea.Cmd
	dbm.DumpPathInput.Focus()
	dbm.DumpPathInput, cmd = dbm.DumpPathInput.Update(msg)
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
		return dbm.updatePathInputView(msg)
	}

	if dbm.Done {
		qh.CreateDump(dbm.SelectedConnectionDetails, dbm.InputDumpPathValue, dbm.SelectedDatabases[0])
		dbm.Done = false
		return dbm, btea.Quit
	}

	return dbm, nil
}

func (dbm DbDumpModel) View() string {

	if dbm.ChoosingConnection {
		return fmt.Sprintf("Select a stored connection: \n%s", dbm.StoredConnectionsList.View())
	}

	if dbm.ChoosingDatabases {
		// Handle selection view
		return fmt.Sprintf("Select a one or more databases: \n%s", dbm.DatabasesList.View())
	}

	if dbm.ChoosingTables {
		// Handle selection view
		return fmt.Sprintf("Select a one or more tables: \n%s", dbm.TablesList.View())
	}

	if dbm.TypingPath {
		return fmt.Sprintf(pathInputTitleStyle.Render("Select a one or more tables:")+" \n\n%s\n\n"+blurredStyle.Render("[ALT+Backspace] go back • [CTRL+C] quit"), dbm.DumpPathInput.View())
	}

	if err := dbm.Err; err != nil {
		return fmt.Sprintf("Sorry, could not fetch tables: \n%s", err)
	}

	return "Press Ctrl+C to Exit"
}

func RunDumpCommand(config *config.SeraphimConfig) {
	seraphimConfig = *config

	input := textinput.New()
	input.Cursor.Style = cursorStyle
	input.TextStyle = focusedStyle
	input.PlaceholderStyle = focusedStyle
	input.Prompt = focusedStyle.Render("\u276F ")

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
		DumpPathInput:         input,
	}

	p := btea.NewProgram(initialModel, btea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		fmt.Printf("FATAL -- Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
