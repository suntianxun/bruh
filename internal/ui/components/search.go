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
	Progress     ProgressModel
	width        int
	height       int
}

func NewSearchModel(width, height int) SearchModel {
	ti := textinput.New()
	ti.Placeholder = "Search remote packages (Enter to search)..."
	ti.Focus()
	ti.CharLimit = 156

	prog := NewProgressModel()
	prog.Message = "Searching..."

	// Calculate 2/3 dimensions for initial setup
	boxWidth := (width * 2) / 3
	boxHeight := (height * 2) / 3
	
	ti.Width = boxWidth - 6

	return SearchModel{
		Input:        ti,
		Results:      NewPackageTable(nil, boxWidth-6, boxHeight-7),
		InputFocused: true,
		Progress:     prog,
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

	m.Progress, cmd = m.Progress.Update(msg)
	cmds = append(cmds, cmd)

	boxWidth := (m.width * 2) / 3
	boxHeight := (m.height * 2) / 3

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		boxWidth = (m.width * 2) / 3
		boxHeight = (m.height * 2) / 3
		m.Input.Width = boxWidth - 6
		m.Results.SetSize(boxWidth-6, boxHeight-7)
	case tea.KeyMsg:
		if m.InputFocused && !m.Progress.Active {
			switch msg.Type {
			case tea.KeyEnter:
				if m.Input.Value() != "" {
					m.Input.Blur()
					m.InputFocused = false
					m.Progress.Active = true
					return m, tea.Batch(m.Progress.Spinner.Tick, SearchCmd(m.Input.Value()))
				}
			case tea.KeyDown, tea.KeyTab:
				m.Input.Blur()
				m.InputFocused = false
			}
		} else if !m.Progress.Active {
			// List focused
			switch msg.Type {
			case tea.KeyUp, tea.KeyShiftTab:
				if m.Results.Index() == 0 || msg.Type == tea.KeyShiftTab {
					m.InputFocused = true
					m.Input.Focus()
				}
			}
		}
	case SearchResultsMsg:
		m.Progress.Active = false
		m.pkgs = msg
		m.Results = NewPackageTable(msg, boxWidth-6, boxHeight-7)
	case error:
		m.Progress.Active = false
	}

	if m.InputFocused {
		m.Input, cmd = m.Input.Update(msg)
		cmds = append(cmds, cmd)
	} else if !m.Progress.Active {
		m.Results, cmd = m.Results.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m SearchModel) View() string {
	boxWidth := (m.width * 2) / 3
	boxHeight := (m.height * 2) / 3

	titleStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#cba6f7")).
		Foreground(lipgloss.Color("#1e1e2e")).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1)

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#cba6f7")).
		Padding(1, 2).
		Width(boxWidth).
		Height(boxHeight).
		Align(lipgloss.Left)

	title := titleStyle.Render(" Remote Search ")

	inputView := m.Input.View()
	if m.InputFocused {
		inputView = lipgloss.NewStyle().Foreground(lipgloss.Color("#cba6f7")).Render(inputView)
	}

	var innerContent string
	if m.Progress.Active {
		// Center the progress inside the table area
		progView := lipgloss.Place(boxWidth-6, boxHeight-6, lipgloss.Center, lipgloss.Center, m.Progress.View())
		innerContent = lipgloss.JoinVertical(lipgloss.Left, inputView, "", progView)
	} else {
		tableHeader := RenderTableHeader(boxWidth - 6)
		innerContent = lipgloss.JoinVertical(lipgloss.Left,
			inputView,
			"",
			tableHeader,
			m.Results.View(),
		)
	}

	content := lipgloss.JoinVertical(lipgloss.Center, title, boxStyle.Render(innerContent))

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content)
}
