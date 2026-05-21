// internal/ui/model.go
package ui

import (
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/brew-tui/internal/brew"
	"github.com/user/brew-tui/internal/ui/components"
)

type activeTab int

const (
	tabInstalled activeTab = iota
	tabSearch
)

type Model struct {
	tab          activeTab
	pkgTable     table.Model
	pkgs         []brew.PackageInfo
	search       components.SearchModel
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
		pkgTable:     components.NewPackageTable(nil, 80, 10),
		search:       components.NewSearchModel(80, 10),
		progress:     prog,
	}
}

type pkgsLoadedMsg []brew.PackageInfo
type errMsg error

func fetchInstalledCmd() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			res, err := brew.GetInstalled()
			if err != nil {
				return errMsg(err)
			}
			return pkgsLoadedMsg(res)
		},
	)
}

func brewActionCmd(action string, name string, isCask bool) tea.Cmd {
	return func() tea.Msg {
		var err error
		switch action {
		case "upgrade":
			err = brew.Upgrade(name, isCask)
		case "uninstall":
			err = brew.Uninstall(name, isCask)
		case "reinstall":
			err = brew.Reinstall(name, isCask)
		}
		
		if err != nil {
			return errMsg(err)
		}
		
		// Refetch packages after an action
		res, err := brew.GetInstalled()
		if err != nil {
			return errMsg(err)
		}
		return pkgsLoadedMsg(res)
	}
}


func (m Model) Init() tea.Cmd {
	return tea.Batch(components.NewProgressModel().Spinner.Tick, fetchInstalledCmd())
}

func getPkgByName(pkgs []brew.PackageInfo, name string) *brew.PackageInfo {
	for _, p := range pkgs {
		if p.Name == name {
			return &p
		}
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.progress, cmd = m.progress.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.pkgTable = components.NewPackageTable(m.pkgs, m.width, m.height-6)
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
			m.tab = (m.tab + 1) % 2
			return m, nil
		case "shift+tab", "left", "h":
			m.tab = (m.tab - 1 + 2) % 2
			return m, nil
		case "u", "d", "r":
			if m.tab == tabInstalled {
				selected := m.pkgTable.SelectedRow()
				if len(selected) > 0 {
					name := selected[0]
					pkg := getPkgByName(m.pkgs, name)
					if pkg != nil {
						m.progress.Active = true
						action := "upgrade"
						if msg.String() == "d" {
							action = "uninstall"
						} else if msg.String() == "r" {
							action = "reinstall"
						}
						m.progress.Message = action + "ing " + name + "..."
						return m, tea.Batch(m.progress.Spinner.Tick, brewActionCmd(action, name, pkg.IsCask))
					}
				}
			}
		}
	case pkgsLoadedMsg:
		m.progress.Active = false
		m.pkgs = msg
		m.pkgTable = components.NewPackageTable(msg, m.width, m.height-6)
	case errMsg:
		m.progress.Active = false
		m.progress.Message = "Error: " + msg.Error()
	}
	
	if !m.progress.Active {
		if m.tab == tabInstalled {
			m.pkgTable, cmd = m.pkgTable.Update(msg)
			cmds = append(cmds, cmd)
		} else if m.tab == tabSearch {
			m.search, cmd = m.search.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	header := HeaderStyle.Render("Homebrew TUI")
	tabs := []string{"Installed", "Search"}
	renderedTabs := components.RenderTabs(tabs, int(m.tab), m.width)
	
	var mainContent string
	if m.progress.Active {
		mainContent = m.progress.View()
	} else {
		switch m.tab {
		case tabInstalled:
			mainContent = lipgloss.NewStyle().Padding(1, 2).Render(m.pkgTable.View())
		case tabSearch:
			mainContent = m.search.View()
		}
	}

	help := "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Render(" u: upgrade • d: uninstall • r: reinstall • tab: switch view • q: quit")
	if m.tab != tabInstalled || m.progress.Active {
		help = ""
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, renderedTabs, mainContent, help)
}