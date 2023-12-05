package selector

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"seraphim/config"
	"time"

	"github.com/JamesStewy/go-mysqldump"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	pgcommands "github.com/habx/pg-commands"
)

var (
	selected        []string
	conn            config.StoredConnection
	defaultDumpPath string
	selectedDb      string
)

type MultiSelectListModel struct {
	PathInput textinput.Model

	Typing bool
	Done   bool

	Choices  []string
	Cursor   int
	Selected map[int]string
}

type MultiSelectResult struct {
	Err    error
	Result map[int]string
}

func (m MultiSelectListModel) Init() tea.Cmd {
	return textinput.Blink
}

func ClearScreen() {

	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls") // for Windows
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear") // for Linux and MacOS
		cmd.Stdout = os.Stdout
		cmd.Run()

	}
}

func (m MultiSelectListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			ClearScreen()
			os.Exit(0)
		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			}
		case "down":
			if m.Cursor < len(m.Choices)-1 {
				m.Cursor++
			}
		case " ":
			_, ok := m.Selected[m.Cursor]
			if ok {
				delete(m.Selected, m.Cursor)
			} else {
				m.Selected[m.Cursor] = m.Choices[m.Cursor]
			}

		case "enter":
			if m.Typing {
				var path string
				if m.PathInput.Value() != "" {
					path = m.PathInput.Value()
				} else {
					path = defaultDumpPath
				}
				res := CreateDump(conn, path, selectedDb)
				if res {
					ClearScreen()
					os.Exit(0)
				} else {
					fmt.Printf("Something went wrong")
					ClearScreen()
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
	return m, nil
}

func (m MultiSelectListModel) View() string {
	s := "Select the desired tables: \n\n"
	for i, choice := range m.Choices {
		cursor := " " // no cursor
		if m.Cursor == i {
			cursor = "\u27A4" // cursor!
		}

		checked := " " // not selected
		if _, ok := m.Selected[i]; ok {
			checked = "\u2714" // selected!
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}
	qs := "\nPress CTRL+C to quit.\n"
	s += qs
	if m.Typing {
		m.PathInput.Focus()
		return fmt.Sprintf("Enter the path for the database dump:\n\n  %s\n\nPress CTRL+C to quit", m.PathInput.View())
	}
	return s
}

func RunMultiSelectList(tables []string, dbConfig config.StoredConnection, db string, defaultDmpPath string) {

	conn = dbConfig
	selectedDb = db
	defaultDumpPath = defaultDmpPath

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

	dumpDir := dumpPath
	dumpFilenameFormat := fmt.Sprintf("%s-%v", dbname, time.Now().Unix())
	switch driver {
	case "mysql":
		db, err := sql.Open(driver, fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", username, password, hostname, port, dbname))
		if err != nil {
			fmt.Println("Error opening database: ", err)
			return false
		}

		dumper, err := mysqldump.Register(db, dumpDir, dumpFilenameFormat)
		if err != nil {
			fmt.Println("Error registering databse:", err)
			return false
		}

		if _, err := dumper.Dump(); err != nil {
			fmt.Println("Error dumping:", err)
			return false
		}
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
			return true
		}
	default:
		log.Fatal("Unknown driver")
		return false
	}
	return true
}
