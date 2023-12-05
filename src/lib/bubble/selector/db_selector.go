package selector

import (
	"errors"
	"fmt"
	"os"
	"seraphim/lib/config"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type DbSelector struct {
	List         list.Model
	Spinner      spinner.Model
	DelegateKeys *delegateKeyMap
	Err          error

	AvailableDbs              []string
	SelectedDb                string
	SelectedConnectionDetails config.StoredConnection
	Choosing                  bool
	Loading                   bool
}

func (ds DbSelector) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (ds DbSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		ds.List.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return ds, tea.Quit
		case "enter":
			ds.Choosing = false
			ds.Loading = false
			ds.SelectedDb = ds.List.SelectedItem().(listItem).Title()
			return ds, tea.Batch(spinner.Tick, tea.Quit)
		}
	case SelectorResult:
		ds.Loading = false
		ds.Choosing = true
		return ds, tea.Quit
	}

	if ds.Loading {
		var cmd tea.Cmd
		ds.Spinner, cmd = ds.Spinner.Update(msg)
		return ds, cmd
	}

	if ds.Choosing {
		var cmd tea.Cmd
		ds.List, cmd = ds.List.Update(msg)
		return ds, cmd
	}

	return ds, nil
}

func (ds DbSelector) View() string {
	if ds.Choosing {
		return fmt.Sprintf("Select the database: \n%s", ds.List.View())
	}

	if ds.Loading {
		return fmt.Sprintf("%s Creating dump", ds.Spinner.View())
	}

	if err := ds.Err; err != nil {
		return fmt.Sprintf("Sorry, could not fetch tables: \n%s", err)
	}

	return "Press Ctrl+C to Exit"
}

func RunDbSelector(selectedDb config.StoredConnection, dbs []string) SelectorResult {
	numItems := len(dbs)
	items := make([]list.Item, numItems)
	delegateKeys := newDelegateKeyMap()
	for i, value := range dbs {
		items[i] = listItem{
			title:       value,
			description: "",
		}
	}

	delegate := newItemDelegate(delegateKeys)
	StoredConnectionList := list.New(items, delegate, 0, 0)
	StoredConnectionList.SetShowFilter(true)
	StoredConnectionList.SetShowTitle(false)
	StoredConnectionList.Styles.Title = titleStyle

	s := spinner.New()
	s.Spinner = spinner.Dot

	initialModel := DbSelector{
		List:         StoredConnectionList,
		AvailableDbs: dbs,
		SelectedDb:   "",
		Spinner:      s,
		Choosing:     true,
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	model, err := p.Run()
	if err != nil {
		fmt.Printf("FATAL -- Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	selected := model.(DbSelector).SelectedDb
	if selected != "" {
		return SelectorResult{
			Err:    nil,
			Result: []string{selected},
		}
	} else {
		if selectedDb.DefaltDatabase != "" {
			return SelectorResult{
				Err:    nil,
				Result: []string{selectedDb.DefaltDatabase},
			}
		} else {
			return SelectorResult{
				Err: errors.New("no default database and no selected database"),
			}
		}
	}

}
