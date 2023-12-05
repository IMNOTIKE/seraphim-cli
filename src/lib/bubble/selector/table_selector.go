package selector

import (
	"fmt"
	"os"
	"seraphim/lib/config"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type TableSelectorMode struct {
	List         list.Model
	Spinner      spinner.Model
	DelegateKeys *delegateKeyMap
	Err          error

	AvailableTables           []string
	SelectedTables            []string
	SelectedConnectionDetails config.StoredConnection
	Choosing                  bool
	Loading                   bool
	PathInput                 textinput.Model
}

func (tbm TableSelectorMode) Init() tea.Cmd {
	return tea.EnterAltScreen
}

type SelectSuccessMsg struct {
	Err   error
	Saved bool
}

func (tbm TableSelectorMode) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		tbm.List.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return tbm, tea.Quit
		case "enter":
			tbm.Choosing = false
			tbm.Loading = true
			tbm.SelectedTables = []string{tbm.List.SelectedItem().(listItem).Title()}
			return tbm, tea.Batch(tbm.CreateDump(), spinner.Tick)
		}
	}
	if tbm.Loading {
		var cmd tea.Cmd
		tbm.Spinner, cmd = tbm.Spinner.Update(msg)
		return tbm, cmd
	}

	if tbm.Choosing {
		var cmd tea.Cmd
		tbm.List, cmd = tbm.List.Update(msg)
		return tbm, cmd
	}

	return tbm, nil
}

func (tbm TableSelectorMode) View() string {

	if tbm.Choosing {
		return fmt.Sprintf("Select one ore more tables: \n%s", tbm.List.View())
	}

	if tbm.Loading {
		return fmt.Sprintf("%s Creating dump", tbm.Spinner.View())
	}

	if err := tbm.Err; err != nil {
		return fmt.Sprintf("Sorry, could not fetch tables: \n%s", err)
	}

	return "Press Ctrl+C to Exit"
}

func (tbm TableSelectorMode) CreateDump() tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)
		tbm.Choosing = false
		tbm.Loading = false
		return SelectSuccessMsg{
			Err:   nil,
			Saved: false,
		}
	}
}

// TODO: Should allow for multi-select
func RunTableSelector(db string, selectedDb config.StoredConnection, tables []string) []string {
	numItems := len(tables)
	items := make([]list.Item, numItems)
	delegateKeys := newDelegateKeyMap()
	for i, value := range tables {
		items[i] = listItem{
			title:       value,
			description: db,
		}
	}

	delegate := newItemDelegate(delegateKeys)
	StoredConnectionList := list.New(items, delegate, 0, 0)
	StoredConnectionList.SetShowFilter(true)
	StoredConnectionList.SetShowTitle(false)
	StoredConnectionList.Styles.Title = titleStyle

	s := spinner.New()
	s.Spinner = spinner.Dot

	initialModel := TableSelectorMode{
		List:            StoredConnectionList,
		AvailableTables: tables,
		SelectedTables:  nil,
		Spinner:         s,
		Choosing:        true,
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	model, err := p.Run()
	if err != nil {
		fmt.Printf("FATAL -- Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	selectedTables := model.(TableSelectorMode).SelectedTables
	if selectedTables != nil {
		return selectedTables
	} else {
		return nil
	}

}
