package util

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FFFF00")).
		Bold(true).
		Padding(0, 1)
)

//------------------------------------------------//
//
//           List items declarations
//
//------------------------------------------------//

type ConnListItem struct {
	Tag  string
	Host string
	User string
}

func (c ConnListItem) Title() string       { return c.Tag }
func (c ConnListItem) Description() string { return c.User + "@" + c.Host }
func (c ConnListItem) FilterValue() string { return c.Tag }

type DbListItem struct {
	Name     string
	Selected bool
}

func (d DbListItem) Title() string { return d.Name }
func (d DbListItem) Description() string {
	if d.Selected {
		return titleStyle.Render("\u2713")
	}
	return ""
}
func (d DbListItem) FilterValue() string { return d.Name }

type TableListItem struct {
	Name     string
	Db       string
	Selected bool
}

func (t TableListItem) Title() string { return t.Name }
func (t TableListItem) Description() string {
	if t.Selected {
		return t.Db + " | " + titleStyle.Render("\u2713")
	}
	return t.Db
}
func (t TableListItem) FilterValue() string { return t.Name }
