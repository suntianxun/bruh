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
	Input        textinput.Model
	Results      list.Model
	pkgs         []brew.PackageInfo
	InputFocused bool
	width        int
	height       int
}

func NewSearchModel(width, height int) SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Search remote packages (Enter to search)..."
	ti.Focus()
	ti.CharLimit = 156
	ti.Width = width - 4

	return SearchModel{
		Input:        ti,
		Results:      NewPackageTable(nil, width-4, height-7),
		InputFocused: true,
		width:        width,
		height:       height,
	}
}

type SearchResultsMsg []brew.PackageInfo

func SearchCmd(query string) tea.Cmd {
	return func() tea.Msg {
		res, err := brew.SearchRemote(query)
		if err != nil {
			return err
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
		m.Input.Width = m.width - 4
		m.Results.SetSize(m.width-4, m.height-7)
	case tea.KeyMsg:
		if m.InputFocused {
			switch msg.Type {
			case tea.KeyEnter:
				if m.Input.Value() != "" {
					m.Input.Blur()
					m.InputFocused = false
					return m, SearchCmd(m.Input.Value())
				}
			case tea.KeyDown, tea.KeyTab:
				m.Input.Blur()
				m.InputFocused = false
			}
		} else {
			// List focused
			switch msg.Type {
			case tea.KeyUp, tea.KeyShiftTab:
				if m.Results.Index() == 0 || msg.Type == tea.KeyShiftTab {
					m.InputFocused = true
					m.Input.Focus()
				}
			}
			// Let list handle navigation
		}
	case SearchResultsMsg:
		m.pkgs = msg
		m.Results = NewPackageTable(msg, m.width-4, m.height-7)
	}

	if m.InputFocused {
		m.Input, cmd = m.Input.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.Results, cmd = m.Results.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m SearchModel) View() string {
	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#cba6f7")). // Mauve
		Padding(1, 1).
		Width(m.width - 2).
		Height(m.height - 2)

	tableHeader := RenderTableHeader(m.width - 4)
	
	inputView := m.Input.View()
	if m.InputFocused {
		inputView = lipgloss.NewStyle().Foreground(lipgloss.Color("#cba6f7")).Render(inputView)
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		inputView,
		"",
		tableHeader,
		m.Results.View(),
	)

	return boxStyle.Render(content)
}