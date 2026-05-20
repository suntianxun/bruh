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
		m.pkgList.SetSize(msg.Width, msg.Height-4) // account for header/footer
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	case pkgsLoadedMsg:
		m.pkgList = components.NewList(msg, m.width, m.height-4)
	}
	
	m.pkgList, cmd = m.pkgList.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	header := HeaderStyle.Render("Homebrew TUI")
	return lipgloss.JoinVertical(lipgloss.Left, header, m.pkgList.View())
}