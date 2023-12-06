package db

import (
	"errors"
	"fmt"
	"seraphim/lib/config"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type adcFormModel struct {
	focusIndex int
	fields     []textinput.Model
	cursorMode cursor.Mode
}

type AdcResult struct {
	Tag           string
	NewConnection config.StoredConnection
	Err           error
}

var (
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))

	newConnTag      string
	newConnHost     string
	newConnUser     string
	newConnPwd      string
	newConnPort     int
	newConnProvider string
	newConnDefDb    string
)

func initialModel() adcFormModel {
	fm := adcFormModel{
		fields: make([]textinput.Model, 7),
	}

	var t textinput.Model
	for i := range fm.fields {
		t = textinput.New()
		t.Cursor.Style = cursorStyle
		t.CharLimit = 56

		switch i {
		case 0:
			t.Placeholder = "Tag"
			t.Focus()
			t.PromptStyle = focusedStyle
			t.TextStyle = focusedStyle
		case 1:
			t.Placeholder = "Host"
			t.CharLimit = 64
		case 2:
			t.Placeholder = "User"
			t.CharLimit = 64
		case 3:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		case 4:
			t.Placeholder = "Port"
			t.CharLimit = 4
		case 5:
			t.Placeholder = "Provider"
			t.CharLimit = 64
		case 6:
			t.Placeholder = "Default database"
			t.CharLimit = 64
		}

		fm.fields[i] = t
	}

	return fm
}

func (fm adcFormModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

func (fm adcFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return fm, tea.Quit

		// Change cursor mode
		case "ctrl+r":
			fm.cursorMode++
			if fm.cursorMode > cursor.CursorHide {
				fm.cursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(fm.fields))
			for i := range fm.fields {
				cmds[i] = fm.fields[i].Cursor.SetMode(fm.cursorMode)
			}
			return fm, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" {
				if fm.focusIndex == len(fm.fields) {
					return fm, tea.Quit
				}
				switch fm.focusIndex {
				case 0:
					newConnTag = fm.fields[fm.focusIndex].Value()
				case 1:
					newConnHost = fm.fields[fm.focusIndex].Value()
				case 2:
					newConnUser = fm.fields[fm.focusIndex].Value()
				case 3:
					newConnPwd = fm.fields[fm.focusIndex].Value()
				case 4:
					if p, err := strconv.Atoi(fm.fields[fm.focusIndex].Value()); err == nil {
						newConnPort = p
					} else {
						return fm, tea.Quit
					}
				case 5:
					newConnProvider = fm.fields[fm.focusIndex].Value()
				case 6:
					newConnDefDb = strings.Replace(strings.Trim(fm.fields[fm.focusIndex].Value(), " "), " ", "_", -1)
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				fm.focusIndex--
			} else {
				fm.focusIndex++
			}

			if fm.focusIndex > len(fm.fields) {
				fm.focusIndex = 0
			} else if fm.focusIndex < 0 {
				fm.focusIndex = len(fm.fields)
			}

			cmds := make([]tea.Cmd, len(fm.fields))
			for i := 0; i <= len(fm.fields)-1; i++ {

				if i == fm.focusIndex {
					// Set focused state
					cmds[i] = fm.fields[i].Focus()
					fm.fields[i].PromptStyle = focusedStyle
					fm.fields[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				fm.fields[i].Blur()
				fm.fields[i].PromptStyle = noStyle
				fm.fields[i].TextStyle = noStyle
			}

			return fm, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := fm.updatefields(msg)

	return fm, cmd
}

func (fm *adcFormModel) updatefields(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(fm.fields))

	// Only text fields with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range fm.fields {
		fm.fields[i], cmds[i] = fm.fields[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (fm adcFormModel) View() string {
	var b strings.Builder

	for i := range fm.fields {
		b.WriteString(fm.fields[i].View())
		if i < len(fm.fields)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if fm.focusIndex == len(fm.fields) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(fm.cursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style) or (ctrl+c or esc to quit)"))

	return b.String()
}

func RunAdcForm() AdcResult {

	var r AdcResult

	r.Err = nil

	if _, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run(); err != nil {
		r.Err = errors.New("Err: " + err.Error())
	}

	r.Tag = newConnTag
	r.NewConnection = config.StoredConnection{
		Host:           newConnHost,
		User:           newConnUser,
		Port:           newConnPort,
		Password:       newConnPwd,
		Provider:       newConnProvider,
		DefaltDatabase: newConnDefDb,
	}

	return r

}
