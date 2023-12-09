package db

import (
	"fmt"
	"os"
	"seraphim/lib/config"
	"seraphim/lib/util"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	appConfig = config.SeraphimConfig{}
)

func (m StoredConnectionEditorModel) updateEditingView(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "alt+backspace", "esc":
			m.Choosing = true
			m.Editing = false
			var cmd tea.Cmd
			m.StoredConnectionsList, cmd = m.StoredConnectionsList.Update(msg)
			return m, tea.Batch(tea.ClearScreen, cmd)
		case "ctrl+r":
			m.CursorMode++
			if m.CursorMode > cursor.CursorHide {
				m.CursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(m.Fields))
			for i := range m.Fields {
				cmds[i] = m.Fields[i].Cursor.SetMode(m.CursorMode)
			}
			return m, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" {
				if m.FocusIndex == len(m.Fields) {
					m.Completed = true
					newConnTag = m.Fields[0].Value()
					if newConnTag == "" {
						newConnTag = m.ChosenConnectionTag
					}
					newConnHost = m.Fields[1].Value()
					if newConnHost == "" {
						newConnHost = m.ChosenConnection.Host
					}
					newConnUser = m.Fields[2].Value()
					if newConnUser == "" {
						newConnUser = m.ChosenConnection.User
					}
					newConnPwd = m.Fields[3].Value()
					if newConnPwd == "" {
						newConnPwd = m.ChosenConnection.Password
					}
					var port = 0
					portFieldValue := m.Fields[4].Value()
					if portFieldValue != "" {
						if p, err := strconv.Atoi(portFieldValue); err == nil {
							port = p
						}
					}
					newConnPort = port
					if newConnPort < 1024 && newConnPort > 49151 {
						newConnPort = m.ChosenConnection.Port
					}
					newConnProvider = m.Fields[5].Value()
					if newConnProvider == "" {
						newConnProvider = m.ChosenConnection.Provider
					}
					newConnDefDb = strings.Replace(strings.Trim(m.Fields[6].Value(), " "), " ", "_", -1)
					if newConnDefDb == "" {
						newConnDefDb = m.ChosenConnection.DefaultDatabase
					}
					newConn := config.StoredConnection{
						Host:            newConnHost,
						User:            newConnUser,
						Password:        newConnPwd,
						Port:            newConnPort,
						Provider:        newConnProvider,
						DefaultDatabase: newConnDefDb,
					}
					m.EditResult = config.EditConnection(appConfig, m.ChosenConnection, newConn, m.ChosenConnectionTag, newConnTag)

					return m, tea.Quit
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				m.FocusIndex--
			} else {
				m.FocusIndex++
			}

			if m.FocusIndex > len(m.Fields) {
				m.FocusIndex = 0
			} else if m.FocusIndex < 0 {
				m.FocusIndex = len(m.Fields)
			}

			cmds := make([]tea.Cmd, len(m.Fields))
			for i := 0; i <= len(m.Fields)-1; i++ {

				if i == m.FocusIndex {
					// Set focused state
					cmds[i] = m.Fields[i].Focus()
					m.Fields[i].PromptStyle = focusedStyle
					m.Fields[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				m.Fields[i].Blur()
				m.Fields[i].PromptStyle = noStyle
				m.Fields[i].TextStyle = noStyle
			}

			return m, tea.Batch(cmds...)
		}
	}

	// Handle character input and blinking
	cmd := m.updateFields(msg)

	return m, cmd
}

func (fm *StoredConnectionEditorModel) updateFields(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(fm.Fields))

	// Only text Fields with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range fm.Fields {
		fm.Fields[i], cmds[i] = fm.Fields[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (fm StoredConnectionEditorModel) viewEditingForm() string {
	var b strings.Builder

	for i := range fm.Fields {
		b.WriteString(fm.Fields[i].View())
		if i < len(fm.Fields)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if fm.FocusIndex == len(fm.Fields) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(fm.CursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style) or (ctrl+c or esc to quit)"))

	return b.String()
}

func (m *StoredConnectionEditorModel) getEditableFields(storedConnection config.StoredConnection, tag string) {

	m.Fields = make([]textinput.Model, 0)

	tagInput := textinput.New()
	tagInput.Placeholder = "Tag: \u21BA " + tag
	tagInput.Focus()
	tagInput.PromptStyle = focusedStyle
	tagInput.TextStyle = focusedStyle
	m.Fields = append(m.Fields, tagInput)
	hostInput := textinput.New()
	hostInput.Placeholder = "Host: \u21BA " + storedConnection.Host
	hostInput.CharLimit = 64
	m.Fields = append(m.Fields, hostInput)
	userInput := textinput.New()
	userInput.Placeholder = "User: \u21BA " + storedConnection.User
	userInput.CharLimit = 64
	m.Fields = append(m.Fields, userInput)
	pwdInput := textinput.New()
	pwdInput.Placeholder = "Password: \u21BA " + strings.Repeat("•", 12)
	pwdInput.EchoMode = textinput.EchoPassword
	pwdInput.EchoCharacter = '•'
	m.Fields = append(m.Fields, pwdInput)
	portInput := textinput.New()
	portInput.Placeholder = fmt.Sprintf("Port: \u21BA %d", storedConnection.Port)
	portInput.CharLimit = 4
	m.Fields = append(m.Fields, portInput)
	providerInput := textinput.New()
	providerInput.Placeholder = "Provider: \u21BA " + storedConnection.Provider
	providerInput.CharLimit = 64
	m.Fields = append(m.Fields, providerInput)
	defDbInput := textinput.New()
	defDbInput.Placeholder = "Default db: \u21BA " + storedConnection.DefaultDatabase
	defDbInput.CharLimit = 64
	m.Fields = append(m.Fields, defDbInput)

}

func (m StoredConnectionEditorModel) updateConnectionChoosingView(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := appStyle.GetFrameSize()
		m.StoredConnectionsList.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			selectedItem := m.StoredConnectionsList.SelectedItem()
			for _, v := range appConfig.Stored_Connections {
				for tag, c := range v {
					if casted, ok := selectedItem.(util.ConnListItem); ok && casted.Tag == tag {
						m.Choosing = false
						m.Editing = true
						m.ChosenConnection = c
						m.ChosenConnectionTag = tag
						m.getEditableFields(c, casted.Tag)
						return m, nil
					}
				}
			}
			return m, tea.Quit // Handle error selected not in list (?) although it should not be possible
		}
	}

	var cmd tea.Cmd
	m.StoredConnectionsList, cmd = m.StoredConnectionsList.Update(msg)
	return m, cmd
}

func RunStoredConnectionEditHandler(sconf *config.SeraphimConfig) {

	appConfig = *sconf
	numItems := len(appConfig.Stored_Connections)
	delegateKeys := newDelegateKeyMap()
	items := make([]list.Item, numItems)
	var i int
	for _, m := range appConfig.Stored_Connections {
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

	editorModel := StoredConnectionEditorModel{
		StoredConnectionsList: StoredConnectionList,
		Choosing:              true,
		ChosenConnectionTag:   "",
	}

	m, err := tea.NewProgram(editorModel, tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if model, ok := m.(StoredConnectionEditorModel); ok && model.EditResult != (config.ConfigOperationResult{}) {
		var indicator string
		if model.EditResult.Err == nil {
			indicator = "\u2713"
		} else {
			indicator = "\u2B59"
		}
		fmt.Println(focusedStyle.Render(fmt.Sprintf("--%s> %s!", indicator, model.EditResult.Msg)))
	}
}

type StoredConnectionEditorModel struct {
	StoredConnectionsList list.Model
	ChosenConnection      config.StoredConnection
	ChosenConnectionTag   string

	Choosing bool
	Editing  bool
	Done     bool

	FocusIndex int
	Fields     []textinput.Model
	CursorMode cursor.Mode
	Completed  bool

	EditResult config.ConfigOperationResult
}

func (m StoredConnectionEditorModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

func (m StoredConnectionEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.Choosing {
		return m.updateConnectionChoosingView(msg)
	}

	if m.Editing {
		return m.updateEditingView(msg)
	}

	return m, nil
}

func (m StoredConnectionEditorModel) View() string {
	if m.Choosing {
		return fmt.Sprintf("Select a stored connection: \n%s", m.StoredConnectionsList.View())
	}

	if m.Editing {
		return m.viewEditingForm()
	}

	return "Press Ctrl+C to exit"
}
