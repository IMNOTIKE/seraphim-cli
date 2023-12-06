package db

import (
	"fmt"
	"log"
	"os"
	"seraphim/lib/config"
	dh "seraphim/lib/db/query"
	"seraphim/lib/util"

	"github.com/charmbracelet/bubbles/list"
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
	DumpPathInput         textinput.Model
	Err                   error

	SelectedConnectionDetails config.StoredConnection
	SelectedDatabases         []util.DbListItem
	SelectedTables            []util.TableListItem
	InputDumpPathValue        string
	ChoosingConnection        bool
	ChoosingDatabases         bool
	ChoosingTables            bool
	TypingPath                bool
	Success                   bool
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

	pathInputTitleStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("230"))
	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	focusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	cursorStyle         = focusedStyle.Copy()

	allTablesSelected = make(map[string]bool, 0)
	anyDbSelected     bool
	anyTableSelcted   bool
)

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
			selectedItem := dbm.StoredConnectionsList.SelectedItem().(util.ConnListItem)
			for _, conn := range seraphimConfig.Stored_Connections {
				t := conn[selectedItem.Tag]
				if t != (config.StoredConnection{}) {
					dbm.SelectedConnectionDetails = t
					break
				}
			}
			dbs := dh.FetchDbList(dbm.SelectedConnectionDetails)
			dbsListItems := make([]list.Item, len(dbs))
			for i, db := range dbs {
				dbsListItems[i] = util.DbListItem{
					Name: db,
				}
				allTablesSelected[db] = false
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
		case "alt+backspace":
			dbm.ChoosingConnection = true
			dbm.SelectedDatabases = make([]util.DbListItem, 0)
			anyDbSelected = false
			dbm.ChoosingDatabases = false
		case " ":
			selectedItem := dbm.DatabasesList.SelectedItem().(util.DbListItem)
			availableItems := dbm.DatabasesList.Items()
			for i, item := range availableItems {
				casted := item.(util.DbListItem)
				if selectedItem == casted {
					if casted.Selected {
						casted.Selected = false

					} else {
						casted.Selected = true
					}
					availableItems[i] = casted
					dbm.DatabasesList.SetItems(availableItems)
				}
			}
			aDs := false
			for _, v := range availableItems {
				casted := v.(util.DbListItem)
				if casted.Selected {
					aDs = true
				}
			}
			anyDbSelected = aDs
		case "enter":
			if anyDbSelected {
				dbm.ChoosingDatabases = false

				selectedDbs := make([]util.DbListItem, 0)
				for _, v := range dbm.DatabasesList.Items() {
					casted := v.(util.DbListItem)
					if casted.Selected {
						selectedDbs = append(selectedDbs, casted)
					}
				}
				dbm.SelectedDatabases = selectedDbs

				tableListItems := make([]list.Item, 0)
				for _, db := range dbm.SelectedDatabases {
					dbTables := dh.FetchTablesForDb(db.Name, dbm.SelectedConnectionDetails)
					tableListItems = append(tableListItems, util.TableListItem{
						Name: "All",
						Db:   db.Name,
					})
					for _, table := range dbTables {
						tableListItems = append(tableListItems, util.TableListItem{
							Name: table,
							Db:   db.Name,
						})
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
			selectedItem := dbm.TablesList.SelectedItem().(util.TableListItem)
			availableItems := dbm.TablesList.Items()
			for i, item := range availableItems {
				casted := item.(util.TableListItem)
				if selectedItem == casted {
					if casted.Selected {
						if casted.Name == "All" {
							allTablesSelected[casted.Db] = false
						}
						casted.Selected = false
					} else {
						if casted.Name == "All" {
							allTablesSelected[casted.Db] = true
							casted.Selected = true
							for k, v := range availableItems {
								vCasted := v.(util.TableListItem)
								if vCasted.Db == casted.Db {
									vCasted.Selected = false
									availableItems[k] = vCasted
								}
							}
						} else if casted.Name != "All" && !allTablesSelected[casted.Db] {
							casted.Selected = true
						}

					}
					availableItems[i] = casted
					dbm.TablesList.SetItems(availableItems)
				}
			}
			aTs := false
			for _, v := range availableItems {
				casted := v.(util.TableListItem)
				if casted.Selected {
					aTs = true
				}
			}
			anyTableSelcted = aTs
		case "alt+backspace":
			dbm.ChoosingDatabases = true
			dbm.SelectedTables = make([]util.TableListItem, 0)
			anyTableSelcted = false
			dbm.ChoosingTables = false
		case "enter":
			if anyTableSelcted {
				dbm.ChoosingTables = false
				for _, v := range dbm.TablesList.Items() {
					casted := v.(util.TableListItem)
					if casted.Selected {
						dbm.SelectedTables = append(dbm.SelectedTables, casted)
					}
				}
				dbm.TypingPath = true
			}
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
			dbm.DumpPathInput.SetValue("")
			dbm.TypingPath = false
		case "enter":
			dbm.TypingPath = false
			inputPath := dbm.DumpPathInput.Value()
			if inputPath == "" {
				inputPath = seraphimConfig.Default_dump_path
			}
			dbm.InputDumpPathValue = inputPath
			dbm.Done = true
			return dbm, btea.Quit
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

	return dbm, nil

}

func (dbm DbDumpModel) View() string {

	s := "Press Ctrl+C to Exit"

	if dbm.ChoosingConnection {
		s = fmt.Sprintf("Select a stored connection: \n%s", dbm.StoredConnectionsList.View())
	}

	if dbm.ChoosingDatabases {
		s = fmt.Sprintf("Select a one or more databases: \n%s", dbm.DatabasesList.View())
	}

	if dbm.ChoosingTables {
		s = fmt.Sprintf("Select a one or more tables: \n%s", dbm.TablesList.View())
	}

	if dbm.TypingPath {
		s = fmt.Sprintf(pathInputTitleStyle.Render("Select a one or more tables:")+" \n\n%s\n\n"+blurredStyle.Render("[ALT+Backspace] go back • [CTRL+C] quit"), dbm.DumpPathInput.View())
	}

	if err := dbm.Err; err != nil {
		return fmt.Sprintf("Sorry, could not fetch tables: \n%s", err)
	}

	return s
}

func RunDumpCommand(sconfig *config.SeraphimConfig) {
	seraphimConfig = *sconfig

	input := textinput.New()
	input.Cursor.Style = cursorStyle
	input.TextStyle = focusedStyle
	input.PlaceholderStyle = focusedStyle
	input.Prompt = focusedStyle.Render("\u276F ")

	numItems := len(sconfig.Stored_Connections)
	delegateKeys := newDelegateKeyMap()
	items := make([]list.Item, numItems)
	var i int
	for _, m := range sconfig.Stored_Connections {
		for key, value := range m {
			items[i] = util.ConnListItem{
				Tag:  key,
				Host: value.Host,
				User: value.User,
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

	initialModel := DbDumpModel{
		StoredConnectionsList: StoredConnectionList,
		ChoosingConnection:    true,
		DumpPathInput:         input,
	}

	p := btea.NewProgram(initialModel, btea.WithAltScreen())
	dbm, err := p.Run()
	if err != nil {
		fmt.Printf("FATAL -- Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	if model, ok := dbm.(DbDumpModel); ok {
		if model.SelectedConnectionDetails != (config.StoredConnection{}) && model.InputDumpPathValue != "" && len(model.SelectedDatabases) != 0 && len(model.SelectedTables) != 0 {
			success := dh.CreateDump(model.SelectedConnectionDetails, model.InputDumpPathValue, model.SelectedDatabases, model.SelectedTables)
			if success {
				fmt.Println(focusedStyle.Render("---> Dump created successfully!"))
			} else {
				fmt.Println(focusedStyle.Render("---> Dump was not created!"))
			}
		}
	} else {
		log.Fatal("something went wrong")
	}

}
