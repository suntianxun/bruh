// internal/ui/model.go
package ui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/brew-tui/internal/brew"
	"github.com/user/brew-tui/internal/ui/components"
)

type activeTab int

const (
	tabInstalled activeTab = iota
	tabSearch
	tabUpdates
)

type Model struct {
	tab      activeTab
	pkgList  list.Model
	width    int
	height   int
}

func InitialModel() Model {
	// Start with empty list, fetch in Init
	return Model{
		tab:     tabInstalled,
		pkgList: components.NewList(nil, 20, 10),
	}
}

type pkgsLoadedMsg []brew.PackageInfo
type errMsg error

func fetchInstalledCmd() tea.Cmd {
	return func() tea.Msg {
		res, err := brew.GetInstalled()
		if err != nil {
			return errMsg(err)
		}
		return pkgsLoadedMsg(res)
	}
}

func (m Model) Init() tea.Cmd {
	return fetchInstalledCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.pkgList.SetSize(msg.Width, m.height-6) // account for header, tabs, footer
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "right", "l":
			m.tab = (m.tab + 1) % 3
			return m, nil
		case "shift+tab", "left", "h":
			m.tab = (m.tab - 1 + 3) % 3
			return m, nil
		}
	case pkgsLoadedMsg:
		m.pkgList = components.NewList(msg, m.width, m.height-6)
	}
	
	if m.tab == tabInstalled {
		m.pkgList, cmd = m.pkgList.Update(msg)
	}
	return m, cmd
}

func (m Model) View() string {
	header := HeaderStyle.Render("Homebrew TUI")
	tabs := []string{"Installed", "Search", "Updates"}
	renderedTabs := components.RenderTabs(tabs, int(m.tab), m.width)
	
	var mainContent string
	switch m.tab {
	case tabInstalled:
		mainContent = m.pkgList.View()
	case tabSearch:
		mainContent = "Search functionality coming soon..."
	case tabUpdates:
		mainContent = "Updates functionality coming soon..."
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, renderedTabs, mainContent)
}