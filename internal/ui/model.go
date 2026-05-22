// internal/ui/model.go
package ui

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
	"github.com/user/brew-tui/internal/brew"
	"github.com/user/brew-tui/internal/ui/components"
)

type logMsg string
type actionDoneMsg struct{ err error }

func waitForLog(c chan string) tea.Cmd {
	return func() tea.Msg {
		msg, ok := <-c
		if !ok {
			return nil
		}
		return logMsg(msg)
	}
}

type Model struct {
	pkgTable   list.Model
	pkgs       []brew.PackageInfo
	search     components.SearchModel
	showSearch bool
	progress   components.ProgressModel
	logChan    chan string
	width      int
	height     int

	animatingHeader bool
	spring          harmonica.Spring
	headerPos       float64
	headerVel       float64
	ticks           int
}

func InitialModel() Model {
	prog := components.NewProgressModel()
	prog.Active = true
	prog.Message = "Loading packages..."
	return Model{
		pkgTable:        components.NewPackageTable(nil, 80, 10, false),
		search:          components.NewSearchModel(80, 10),
		progress:        prog,
		animatingHeader: true,
		spring:          harmonica.NewSpring(harmonica.FPS(60), 4.0, 1.0), // Critically damped, no bouncing
		headerPos:       40.0,
		headerVel:       0.0,
		ticks:           0,
	}
}

type pkgsLoadedMsg []brew.PackageInfo
type errMsg error
type frameMsg time.Time

func animateCmd() tea.Cmd {
	return tea.Tick(time.Second/60, func(t time.Time) tea.Msg {
		return frameMsg(t)
	})
}

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

func brewActionCmd(action string, name string, isCask bool, c chan string) tea.Cmd {
	return func() tea.Msg {
		var err error
		switch action {
		case "upgrade":
			err = brew.Upgrade(name, isCask, c)
		case "uninstall":
			err = brew.Uninstall(name, isCask, c)
		case "reinstall":
			err = brew.Reinstall(name, isCask, c)
		case "install":
			err = brew.Install(name, isCask, c)
		}
		
		close(c)
		return actionDoneMsg{err}
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(components.NewProgressModel().Spinner.Tick, fetchInstalledCmd(), animateCmd())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.progress, cmd = m.progress.Update(msg)
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case frameMsg:
		m.ticks++
		if m.animatingHeader {
			m.headerPos, m.headerVel = m.spring.Update(m.headerPos, m.headerVel, 0)
			if math.Abs(m.headerPos) < 0.1 && math.Abs(m.headerVel) < 0.1 {
				m.headerPos = 0
				m.animatingHeader = false
			}
		}
		cmds = append(cmds, animateCmd())

	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.pkgTable = components.NewPackageTable(m.pkgs, m.width, m.height-16, false)
		m.search, _ = m.search.Update(msg)

	case logMsg:
		m.progress.Logs = append(m.progress.Logs, string(msg))
		if len(m.progress.Logs) > 10 {
			m.progress.Logs = m.progress.Logs[len(m.progress.Logs)-10:]
		}
		return m, tea.Batch(m.progress.Spinner.Tick, waitForLog(m.logChan))

	case actionDoneMsg:
		if msg.err != nil { // an error occurred
			m.progress.Message = "Error: " + msg.err.Error()
			m.progress.Active = false
		} else {
			m.progress.Message = "Reloading packages..."
			m.progress.Logs = nil // clear logs
			return m, tea.Batch(m.progress.Spinner.Tick, fetchInstalledCmd())
		}

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
					selected := m.search.Results.SelectedItem()
					if selected != nil {
						item := selected.(components.PackageItem)
						m.progress.Active = true
						m.progress.Logs = nil
						m.logChan = make(chan string, 100)
						
						action := "install"
						if msg.String() == "d" {
							action = "uninstall"
						}
						m.progress.Message = action + "ing " + item.Name + "..."
						m.showSearch = false // close search on action
						
						return m, tea.Batch(m.progress.Spinner.Tick, waitForLog(m.logChan), brewActionCmd(action, item.Name, item.IsCask, m.logChan))
					}
				}
			}
		} else {
			// Main table view
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
				selected := m.pkgTable.SelectedItem()
				if selected != nil {
					item := selected.(components.PackageItem)
					name := item.Name
					m.progress.Active = true
					m.progress.Logs = nil
					m.logChan = make(chan string, 100)

					action := "upgrade"
					if msg.String() == "d" {
						action = "uninstall"
					} else if msg.String() == "r" {
						action = "reinstall"
					}
					m.progress.Message = action + "ing " + name + "..."
					return m, tea.Batch(m.progress.Spinner.Tick, waitForLog(m.logChan), brewActionCmd(action, name, item.IsCask, m.logChan))
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
		m.pkgTable = components.NewPackageTable(sorted, m.width, m.height-12, false)
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

func RenderGradientBruh(pos float64, ticks int) string {
	ascii := []string{
		`█▀▀▀▄ █▀▀▀▄ █   █ █   █`,
		`█▀▀▀▄ █▀▀▀▄ █   █ █▀▀▀█`,
		`▀▀▀▀  ▀   ▀  ▀▀▀  ▀   ▀`,
	}
	colors := []lipgloss.Color{
		lipgloss.Color("#cba6f7"),
		lipgloss.Color("#f5c2e7"),
		lipgloss.Color("#f2cdcd"),
		lipgloss.Color("#f5e0dc"),
	}

	var out []string
	
	padCount := int(math.Round(pos))
	if padCount < -10 {
		padCount = -10
	}
	padStr := strings.Repeat(" ", 10 + padCount)
	sweep := (ticks % 120)

	for _, line := range ascii {
		var coloredLine string
		for i, c := range line {
			colorIdx := (i * len(colors)) / len(line)
			if colorIdx >= len(colors) {
				colorIdx = len(colors) - 1
			}
			
			fg := colors[colorIdx]
			
			// Color sweep logic
			sweepPos := float64(sweep) / 2.0
			dist := math.Abs(float64(i) - sweepPos)
			if dist < 3.0 {
				if dist < 1.0 {
					fg = lipgloss.Color("#ffffff")
				} else if dist < 2.0 {
					fg = lipgloss.Color("#f9e2af")
				}
			}

			coloredLine += lipgloss.NewStyle().Foreground(fg).Bold(true).Render(string(c))
		}
		out = append(out, padStr + coloredLine)
	}

	rendered := lipgloss.JoinVertical(lipgloss.Left, out...)
	
	// Add padding to the top to distance it from the terminal edge
	rendered = lipgloss.NewStyle().PaddingTop(2).Render(rendered)

	return rendered + "\n\n\n"
}

func (m Model) View() string {
	if m.showSearch {
		searchView := lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, m.search.View())
		return searchView
	}

	header := RenderGradientBruh(m.headerPos, m.ticks)
	
	var mainContent string
	if m.progress.Active {
		mainContent = lipgloss.Place(m.width, m.height-12, lipgloss.Center, lipgloss.Center, m.progress.View())
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