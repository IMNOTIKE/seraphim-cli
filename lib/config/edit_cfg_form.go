package config

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/huh"
)

var (
	appConfig    SeraphimConfig
	focusedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)

func RunCfgEditForm(sconfig *SeraphimConfig) {

	appConfig = *sconfig
	initModel, err := tea.NewProgram(initialModel(), tea.WithAltScreen()).Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if model, ok := initModel.(CfgEditorModel); ok && model.Result != (ConfigOperationResult{}) {
		var indicator string
		var confirm bool
		if model.Result.Err == nil {
			huh.NewConfirm().
			Title("You sure?").
			Affirmative("Yes!").
			Negative("No.").
			Value(&confirm)
			indicator = "\u2713"
		} else {
			indicator = "\u2B59"
		}
		fmt.Println(focusedStyle.Render(fmt.Sprintf("--%s> %s!", indicator, model.Result.Msg)))
	}
}

func (cfgm *CfgEditorModel) getEditableFields() {

	cfgm.Fields = make([]textinput.Model, 0)

	tagInput := textinput.New()
	tagInput.Placeholder = ": \u21BA " + appConfig.Default_dump_path
	tagInput.Focus()
	tagInput.PromptStyle = focusedStyle
	tagInput.TextStyle = focusedStyle
	cfgm.Fields = append(cfgm.Fields, tagInput)

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
	Result ConfigOperationResult
}

func (cfgm CfgEditorModel) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, tea.EnterAltScreen)
}

func (cfgm CfgEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return cfgm, tea.Quit
		}
	}
	return cfgm, nil
}

func (cfgm CfgEditorModel) View() string {
	return "Press Ctrl+C to exit"
}
