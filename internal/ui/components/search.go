// internal/ui/components/search.go
package components

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/brew-tui/internal/brew"
)

type SearchModel struct {
	Input   textinput.Model
	Results list.Model
	pkgs    []brew.PackageInfo
	width   int
	height  int
}

func NewSearchModel(width, height int) SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Search for a package..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = 50

	return SearchModel{
		Input:   ti,
		Results: NewPackageTable(nil, width, height-4),
		width:   width,
		height:  height,
	}
}

type SearchResultsMsg []brew.PackageInfo

func SearchCmd(query string) tea.Cmd {
	return func() tea.Msg {
		res, err := brew.Search(query)
		if err != nil {
			return err // In a real app we'd wrap this
		}
		return SearchResultsMsg(res)
	}
}

func (m SearchModel) Update(msg tea.Msg) (SearchModel, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.Results = NewPackageTable(m.pkgs, m.width, m.height-4)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.Input.Focused() {
				m.Input.Blur()
				return m, SearchCmd(m.Input.Value())
			}
		case tea.KeyEsc:
			m.Input.Focus()
		}
	case SearchResultsMsg:
		m.pkgs = msg
		m.Results = NewPackageTable(msg, m.width, m.height-4)
	}

	if m.Input.Focused() {
		m.Input, cmd = m.Input.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.Results, cmd = m.Results.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m SearchModel) View() string {
	tableHeader := RenderTableHeader(m.width)
	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Padding(1, 2).Render(m.Input.View()),
		tableHeader,
		m.Results.View(),
	)
}