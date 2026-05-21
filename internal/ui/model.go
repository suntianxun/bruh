// internal/ui/model.go
package ui

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/brew-tui/internal/brew"
	"github.com/user/brew-tui/internal/ui/components"
)

type Model struct {
	pkgTable   list.Model
	pkgs       []brew.PackageInfo
	search     components.SearchModel
	showSearch bool
	progress   components.ProgressModel
	width      int
	height     int
}

func InitialModel() Model {
	prog := components.NewProgressModel()
	prog.Active = true
	prog.Message = "Loading packages..."
	return Model{
		pkgTable: components.NewPackageTable(nil, 80, 10),
		search:   components.NewSearchModel(80, 10),
		progress: prog,
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
		case "install":
			err = brew.Install(name, isCask)
		}
		
		if err != nil {
			return errMsg(err)
		}
		
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

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.progress, cmd = m.progress.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.pkgTable = components.NewPackageTable(m.pkgs, m.width, m.height-12)
		m.search, _ = m.search.Update(msg)
	case tea.KeyMsg:
		if m.progress.Active {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			}
			return m, tea.Batch(cmds...)
		}

		if m.showSearch {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "q":
				if !m.search.InputFocused {
					m.showSearch = false
					return m, nil
				}
			case "i", "d":
				if !m.search.InputFocused {
					cursor := m.search.Results.Index()
					if cursor >= 0 && cursor < len(m.search.Results.Items()) {
						item := m.search.Results.Items()[cursor].(components.PackageItem)
						m.progress.Active = true
						action := "install"
						if msg.String() == "d" {
							action = "uninstall"
						}
						m.progress.Message = action + "ing " + item.Name + "..."
						m.showSearch = false // close search on action
						return m, tea.Batch(m.progress.Spinner.Tick, brewActionCmd(action, item.Name, item.IsCask))
					}
				}
			}
		} else {
			// Main table view
			// Let list handle its own filter input if active
			if m.pkgTable.FilterState() == list.Filtering {
				break
			}
			switch msg.String() {
			case "ctrl+c", "q":
				return m, tea.Quit
			case "s":
				m.showSearch = true
				return m, nil
			case "u", "d", "r":
				cursor := m.pkgTable.Index()
				if cursor >= 0 && cursor < len(m.pkgTable.Items()) {
					item := m.pkgTable.Items()[cursor].(components.PackageItem)
					name := item.Name
					m.progress.Active = true
					action := "upgrade"
					if msg.String() == "d" {
						action = "uninstall"
					} else if msg.String() == "r" {
						action = "reinstall"
					}
					m.progress.Message = action + "ing " + name + "..."
					return m, tea.Batch(m.progress.Spinner.Tick, brewActionCmd(action, name, item.IsCask))
				}
			}
		}

	case pkgsLoadedMsg:
		m.progress.Active = false
		sorted := make([]brew.PackageInfo, len(msg))
		copy(sorted, msg)
		sort.Slice(sorted, func(i, j int) bool {
			if sorted[i].IsOutdated == sorted[j].IsOutdated {
				return sorted[i].Name < sorted[j].Name
			}
			return sorted[i].IsOutdated // true (outdated) comes before false
		})
		m.pkgs = sorted
		m.pkgTable = components.NewPackageTable(sorted, m.width, m.height-12)
	case errMsg:
		m.progress.Active = false
		m.progress.Message = "Error: " + msg.Error()
	}
	
	if !m.progress.Active {
		if m.showSearch {
			m.search, cmd = m.search.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			m.pkgTable, cmd = m.pkgTable.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	
	return m, tea.Batch(cmds...)
}

func RenderGradientBruh() string {
	ascii := []string{
		`  ____                  _     `,
		` |  _ \                | |    `,
		` | |_) | _ __  _   _   | |__  `,
		` |  _ < | '__|| | | |  | '_ \ `,
		` | |_) || |   | |_| |  | | | |`,
		` |____/ |_|    \__,_|  |_| |_|`,
	}
	colors := []lipgloss.Color{
		lipgloss.Color("#cba6f7"),
		lipgloss.Color("#f5c2e7"),
		lipgloss.Color("#f2cdcd"),
		lipgloss.Color("#f5e0dc"),
	}

	var out []string
	for _, line := range ascii {
		var coloredLine string
		for i, c := range line {
			colorIdx := (i * len(colors)) / len(line)
			if colorIdx >= len(colors) {
				colorIdx = len(colors) - 1
			}
			coloredLine += lipgloss.NewStyle().Foreground(colors[colorIdx]).Bold(true).Render(string(c))
		}
		out = append(out, coloredLine)
	}
	return lipgloss.JoinVertical(lipgloss.Left, out...) + "\n\n"
}

func (m Model) View() string {
	if m.showSearch {
		// Calculate overlay position manually or just show it centered
		searchView := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.search.View())
		return searchView
	}

	header := RenderGradientBruh()
	
	var mainContent string
	if m.progress.Active {
		mainContent = m.progress.View()
	} else {
		tableHeader := components.RenderTableHeader(m.width)
		mainContent = lipgloss.JoinVertical(lipgloss.Left, tableHeader, m.pkgTable.View())
	}

	var total, outdated, uptodate int
	for _, p := range m.pkgs {
		total++
		if p.IsOutdated {
			outdated++
		} else {
			uptodate++
		}
	}
	
	statusText := fmt.Sprintf(" Total: %d • Up-to-date: %d • Outdated: %d", total, uptodate, outdated)
	statusLine := lipgloss.NewStyle().Foreground(lipgloss.Color("#a6adc8")).Padding(0, 1).Render(statusText)

	helpText := " /: filter • s: search remote • u: upgrade • d: uninstall • r: reinstall • q: quit"
	if m.progress.Active {
		helpText = ""
	}
	help := lipgloss.NewStyle().Foreground(lipgloss.Color("#6c7086")).Render(helpText)

	footer := lipgloss.JoinHorizontal(lipgloss.Left, statusLine, "  ", help)

	return lipgloss.JoinVertical(lipgloss.Left, header, mainContent, footer)
}
