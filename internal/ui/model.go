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
	tab          activeTab
	pkgList      list.Model
	search       components.SearchModel
	outdatedList list.Model
	progress     components.ProgressModel
	width        int
	height       int
}

func InitialModel() Model {
	prog := components.NewProgressModel()
	prog.Active = true
	prog.Message = "Loading packages..."
	return Model{
		tab:          tabInstalled,
		pkgList:      components.NewList(nil, 20, 10),
		search:       components.NewSearchModel(20, 10),
		outdatedList: components.NewList(nil, 20, 10),
		progress:     prog,
	}
}

type pkgsLoadedMsg []brew.PackageInfo
type outdatedLoadedMsg []brew.PackageInfo
type errMsg error

func fetchInstalledCmd() tea.Cmd {
	return tea.Batch(
		components.NewProgressModel().Spinner.Tick,
		func() tea.Msg {
			res, err := brew.GetInstalled()
			if err != nil {
				return errMsg(err)
			}
			return pkgsLoadedMsg(res)
		},
	)
}

func fetchOutdatedCmd() tea.Cmd {
	return tea.Batch(
		components.NewProgressModel().Spinner.Tick,
		func() tea.Msg {
			res, err := brew.GetOutdated()
			if err != nil {
				return errMsg(err)
			}
			return outdatedLoadedMsg(res)
		},
	)
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(fetchInstalledCmd(), fetchOutdatedCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.progress, cmd = m.progress.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.pkgList.SetSize(msg.Width, m.height-6) // account for header, tabs, footer
		m.outdatedList.SetSize(msg.Width, m.height-6)
		m.search, _ = m.search.Update(msg)
	case tea.KeyMsg:
		// Don't allow navigation if progress is active
		if m.progress.Active {
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			}
			return m, tea.Batch(cmds...)
		}

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
		m.progress.Active = false
		m.pkgList = components.NewList(msg, m.width, m.height-6)
	case outdatedLoadedMsg:
		m.progress.Active = false
		m.outdatedList = components.NewList(msg, m.width, m.height-6)
		m.outdatedList.Title = "Updates Available"
	}
	
	if !m.progress.Active {
		if m.tab == tabInstalled {
			m.pkgList, cmd = m.pkgList.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.tab == tabSearch {
			m.search, cmd = m.search.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.tab == tabUpdates {
			m.outdatedList, cmd = m.outdatedList.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	header := HeaderStyle.Render("Homebrew TUI")
	tabs := []string{"Installed", "Search", "Updates"}
	renderedTabs := components.RenderTabs(tabs, int(m.tab), m.width)
	
	var mainContent string
	if m.progress.Active {
		mainContent = m.progress.View()
	} else {
		switch m.tab {
		case tabInstalled:
			mainContent = m.pkgList.View()
		case tabSearch:
			mainContent = m.search.View()
		case tabUpdates:
			mainContent = m.outdatedList.View()
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, renderedTabs, mainContent)
}