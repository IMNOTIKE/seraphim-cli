package db

import (
	"fmt"
	"log"
	"os"
	"seraphim/config"
	"seraphim/lib/bubble/selector"
	qh "seraphim/lib/db/query"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	btea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type DbDumpModel struct {
	List         list.Model
	DelegateKeys *delegateKeyMap
	Spinner      spinner.Model
	DbInput      textinput.Model
	Err          error

	AvailableConnections      []config.StoredConnection
	SelectedConnectionDetails config.StoredConnection
	Database                  string
	Choosing                  bool
	Selecting                 bool
	Tables                    []string
}

type SelectSuccessMsg struct {
	Err       error
	ResultSet []string
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

type listItem struct {
	tag  string
	host string
	user string
}

func (i listItem) Title() string       { return i.tag }
func (i listItem) Description() string { return i.user + "@" + i.host }
func (i listItem) FilterValue() string { return i.tag }

func (dbm DbDumpModel) Update(msg btea.Msg) (btea.Model, btea.Cmd) {
	switch msg := msg.(type) {
	case btea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		dbm.List.SetSize(msg.Width-h, msg.Height-v)
	case btea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return dbm, btea.Quit
		case "enter":
			dbm.Choosing = false
			selectedItem := dbm.List.SelectedItem().(listItem)
			var selectedConn config.StoredConnection
			for _, conn := range seraphimConfig.Stored_Connections {
				t := conn[selectedItem.tag]
				selectedConn = t

			}
			dbm.PerformDump(selectedConn, seraphimConfig.Default_dump_path)
			return dbm, nil
		}
	case SelectSuccessMsg:
		dbm.Choosing = true
		dbm.Tables = msg.ResultSet
		return dbm, btea.Quit
	}

	if dbm.Choosing {
		var cmd btea.Cmd
		dbm.List, cmd = dbm.List.Update(msg)
		return dbm, cmd
	}

	if dbm.Selecting {
		var cmd btea.Cmd
		dbm.DbInput, cmd = dbm.DbInput.Update(msg)
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
			items[i] = listItem{
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
		List:     StoredConnectionList,
		Spinner:  s,
		Choosing: true,
		Tables:   nil,
	}

	p := btea.NewProgram(initialModel, btea.WithAltScreen())
	_, err := p.Run()
	if err != nil {
		fmt.Printf("FATAL -- Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
