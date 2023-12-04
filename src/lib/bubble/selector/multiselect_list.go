package selector

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var selected []string

type MultiSelectListModel struct {
	Choices  []string       // items on the to-do list
	Cursor   int            // which to-do list item our Cursor is pointing at
	Selected map[int]string // which to-do items are Selected
}

func (m MultiSelectListModel) Init() tea.Cmd {
	// Just return `nil`, which means "no I/O right now, please."
	return nil
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
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" and "k" keys move the Cursor up
		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}

		// The "down" and "j" keys move the Cursor down
		case "down", "j":
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
		return m, tea.Quit
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

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

func RunMultiSelectList(tables []string) []string {

	initialModel := MultiSelectListModel{
		Choices:  append([]string{"All"}, tables...),
		Selected: make(map[int]string),
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
	return selected
}
