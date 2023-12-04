package selector

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"seraphim/config"
	"time"

	"github.com/JamesStewy/go-mysqldump"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	pgcommands "github.com/habx/pg-commands"
)

var (
	selected   []string
	conn       config.StoredConnection
	selectedDb string
)

type MultiSelectListModel struct {
	PathInput textinput.Model

	Typing bool
	Done   bool

	Choices  []string       // items on the to-do list
	Cursor   int            // which to-do list item our Cursor is pointing at
	Selected map[int]string // which to-do items are Selected
}

func (m MultiSelectListModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return textinput.Blink
}

type MultiSelectResult struct {
	Err    error
	Result map[int]string
}

func (m MultiSelectListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c":
			return m, tea.Quit

		// The "up" keys move the Cursor up
		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			}

		// The "down" keys move the Cursor down
		case "down":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the Selected state for the item that the Cursor is pointing at.
		case " ":
			_, ok := m.Selected[m.Cursor]
			if ok {
				delete(m.Selected, m.Cursor)
			} else {
				m.Selected[m.Cursor] = m.Choices[m.Cursor]
			}

		case "enter":
			if m.Typing {
				res := CreateDump(conn, m.PathInput.Value(), selectedDb)
				if res {
					os.Exit(0)
				} else {
					fmt.Printf("Something went wrong")
					os.Exit(1)
				}
			}
			return m, func() tea.Msg {
				return MultiSelectResult{
					Err:    nil,
					Result: m.Selected,
				}
			}

		}
	case MultiSelectResult:
		for _, v := range m.Selected {
			selected = append(selected, v)
		}
		m.Typing = true
		return m, nil
	}

	if m.Typing {
		var cmd tea.Cmd
		m.PathInput, cmd = m.PathInput.Update(msg)
		return m, cmd
	}
	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m MultiSelectListModel) View() string {
	// The header
	s := "What should we buy at the market?\n\n"

	// Iterate over our Choices
	for i, choice := range m.Choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.Cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.Selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	qs := "\nPress CTRL+C to quit.\n"
	s += qs
	if m.Typing {
		m.PathInput.Focus()
		return fmt.Sprintf("Enter the path for the database dump:\n\n  %s\n\nPress CTRL+C to quit", m.PathInput.View())
	}
	// The footer

	// Send the UI for rendering
	return s
}

func RunMultiSelectList(tables []string, dbConfig config.StoredConnection, db string) {

	conn = dbConfig
	selectedDb = db

	t := textinput.New()
	t.Focus()
	initialModel := MultiSelectListModel{
		Choices:   append([]string{"All"}, tables...),
		Selected:  make(map[int]string),
		PathInput: t,
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func CreateDump(selected config.StoredConnection, dumpPath string, selectedDb string) bool {
	username := selected.User
	password := selected.Password
	hostname := selected.Host
	port := selected.Port
	dbname := selectedDb
	driver := selected.Provider

	dumpDir := dumpPath                                            // you should create this directory
	dumpFilenameFormat := fmt.Sprintf("%s-%v", dbname, time.Now()) // accepts time layout string and add .sql at the end of file
	// ADD default dump path to config file
	switch driver {
	case "mysql":
		db, err := sql.Open(driver, fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", username, password, hostname, port, dbname))
		if err != nil {
			fmt.Println("Error opening database: ", err)
			return false
		}

		// Register database with mysqldump
		dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
		if err != nil {
			fmt.Println("Error registering databse:", err)
			return false
		}

		// Dump database to file
		resultFilename, err := dumper.Dump()
		if err != nil {
			fmt.Println("Error dumping:", err)
			return false
		}
		fmt.Printf("File is saved to %s", resultFilename)

		// Close dumper and connected database
		dumper.Close()
	case "postgres":
		dump, _ := pgcommands.NewDump(&pgcommands.Postgres{
			Host:     hostname,
			Port:     port,
			DB:       dbname,
			Username: username,
			Password: password,
		})
		dumpExec := dump.Exec(pgcommands.ExecOptions{StreamPrint: false})
		if dumpExec.Error != nil {
			fmt.Println(dumpExec.Error.Err)
			fmt.Println(dumpExec.Output)
			return false

		} else {
			fmt.Println("Dump success")
			fmt.Println(dumpExec.Output)
			return false
		}
	default:
		log.Fatal("Unknown driver")
		return false
	}
	return true
}
