package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	appConfig    SeraphimConfig
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	noStyle             = lipgloss.NewStyle()
	helpStyle           = blurredStyle.Copy()
	cursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	focusedButton = focusedStyle.Copy().Render("[ Submit ]")
	blurredButton = fmt.Sprintf("[ %s ]", blurredStyle.Render("Submit"))

	blurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))

)

type ConfigEditorResult struct {
	Err error `default:"nil"`  
	Msg string
	EditedConfig SeraphimConfig
}

func RunCfgEditForm(sconfig *SeraphimConfig) ConfigEditorResult {

	appConfig = *sconfig
	initModel, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if model, ok := initModel.(CfgEditorModel); ok && model.Result.Err == nil {
		var indicator string
		var confirm bool
		var message string
		if model.Result.Err == nil && model.Completed {
			huh.NewConfirm().
			Title("You sure?").
			Affirmative("Yes!").
			Negative("No.").
			Value(&confirm).Run()
			if confirm {
				if r := SaveConfig(model.Result.EditedConfig); r.Err == nil {
					indicator = "\u2713 "
					message = r.Msg
				}else {
					indicator = "\u2B59 "
					message = r.Err.Error()
				}
			}else {
				indicator = "\u2B59 "
				message = "Aborted operation"
			}
		} else if model.Result.Err != nil {
			message = model.Result.Err.Error()
		}
		fmt.Println(focusedStyle.Render(fmt.Sprintf("--%s> %s!", indicator, message)))
	}
	return ConfigEditorResult{}
}

func (cfgm *CfgEditorModel) getEditableFields() {

	cfgm.Fields = make([]textinput.Model, 0)

	defDumpPathInput := textinput.New()
	defDumpPathInput.Placeholder = "Default dump path: \u21BA " + appConfig.Default_dump_path
	defDumpPathInput.Focus()
	defDumpPathInput.PromptStyle = focusedStyle
	defDumpPathInput.TextStyle = focusedStyle
	cfgm.Fields = append(cfgm.Fields, defDumpPathInput)

}

func initialModel() CfgEditorModel {
	initModel := CfgEditorModel{
		Editing: true,
	}
	initModel.getEditableFields()
	return initModel
}

type CfgEditorModel struct {
	// MultiField form
	FocusIndex int
	Fields     []textinput.Model
	CursorMode cursor.Mode
	// state
	Editing   bool
	Completed bool
	Done      bool
	// result
	Result ConfigEditorResult
}

func (cfgm CfgEditorModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

func (cfgm CfgEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			ClearScreen()
			os.Exit(0)
		case "ctrl+r":
			cfgm.CursorMode++
			if cfgm.CursorMode > cursor.CursorHide {
				cfgm.CursorMode = cursor.CursorBlink
			}
			cmds := make([]tea.Cmd, len(cfgm.Fields))
			for i := range cfgm.Fields {
				cmds[i] = cfgm.Fields[i].Cursor.SetMode(cfgm.CursorMode)
			}
			return cfgm, tea.Batch(cmds...)

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" {
				if cfgm.FocusIndex == len(cfgm.Fields) {
					cfgm.Completed = true
					var confCopy SeraphimConfig
					b, _ := json.Marshal(appConfig)
					json.Unmarshal(b, &confCopy)
					if cfgm.Fields[0].Value() != "" {
						confCopy.Default_dump_path = cfgm.Fields[0].Value() 
					}
					cfgm.Result.EditedConfig = confCopy
					return cfgm, tea.Quit
				}
			}

			// Cycle indexes
			if s == "up" || s == "shift+tab" {
				cfgm.FocusIndex--
			} else {
				cfgm.FocusIndex++
			}

			if cfgm.FocusIndex > len(cfgm.Fields) {
				cfgm.FocusIndex = 0
			} else if cfgm.FocusIndex < 0 {
				cfgm.FocusIndex = len(cfgm.Fields)
			}

			cmds := make([]tea.Cmd, len(cfgm.Fields))
			for i := 0; i <= len(cfgm.Fields)-1; i++ {

				if i == cfgm.FocusIndex {
					// Set focused state
					cmds[i] = cfgm.Fields[i].Focus()
					cfgm.Fields[i].PromptStyle = focusedStyle
					cfgm.Fields[i].TextStyle = focusedStyle
					continue
				}
				// Remove focused state
				cfgm.Fields[i].Blur()
				cfgm.Fields[i].PromptStyle = noStyle
				cfgm.Fields[i].TextStyle = noStyle
			}

			return cfgm, tea.Batch(cmds...)
		}
	}
	cmd := cfgm.updateFields(msg)

	return cfgm, cmd
}

func (cfgm CfgEditorModel) View() string {
	var b strings.Builder

	for i := range cfgm.Fields {
		b.WriteString(cfgm.Fields[i].View())
		if i < len(cfgm.Fields)-1 {
			b.WriteRune('\n')
		}
	}

	button := &blurredButton
	if cfgm.FocusIndex == len(cfgm.Fields) {
		button = &focusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	b.WriteString(helpStyle.Render("cursor mode is "))
	b.WriteString(cursorModeHelpStyle.Render(cfgm.CursorMode.String()))
	b.WriteString(helpStyle.Render(" (ctrl+r to change style) or (ctrl+c / esc to quit)"))

	return b.String()
}

func (cfgm *CfgEditorModel) updateFields(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(cfgm.Fields))

	// Only text Fields with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range cfgm.Fields {
		cfgm.Fields[i], cmds[i] = cfgm.Fields[i].Update(msg)
	}

	return tea.Batch(cmds...)
}
