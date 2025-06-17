package ui

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/crit/fake-ops/internal/app"
)

type QuitSignal struct{}

func DelayedQuit() tea.Msg {
	time.Sleep(1 * time.Second)
	return QuitSignal{}
}

type UI struct {
	ctx      *app.Context
	cancel   func()
	width    int
	height   int
	logs     *LogView
	services *ServiceView
}

func (ui *UI) Init() tea.Cmd {
	return func() tea.Msg {
		// Clear the terminal screen
		fmt.Print("\033[H\033[2J")
		return nil
	}
}

func (ui *UI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Check for the keys "q" or "Ctrl+C" to quit
		switch msg.String() {
		case "q", "ctrl+c":
			ui.cancel()
			return ui, DelayedQuit
		}
		return ui, nil // Ignore all other keys

	case QuitSignal:
		return ui, tea.Quit

	case tea.WindowSizeMsg:
		colw, colh := (msg.Width-4)/2, msg.Height-6

		ui.width = msg.Width
		ui.height = msg.Height
		ui.services.Resize(colw+colw, colh) // TODO: troubleshoot
		ui.logs.Resize(colw, colh)
		return ui, nil

	case app.Message:
		ui.logs.Update(msg)
		return ui, nil

	case app.ServiceMessage:
		ui.services.Update(msg)
		return ui, nil

	case app.ServiceStatus:
		ui.services.UpdateStatus(msg)
		return ui, nil

	default:
		return ui, nil
	}
}

func (ui *UI) View() string {
	var containerStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(cBorder).
		Padding(1).
		Width(ui.width - 2).
		Height(ui.height - 2)

	row1 := lipgloss.JoinHorizontal(lipgloss.Top, ui.services.View(), ui.logs.View())

	footer := lipgloss.NewStyle().
		Foreground(cSecondary).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(cBorder).
		Width(ui.width - 4).
		Height(1).
		Align(lipgloss.Right).
		Render("q or ctrl+c to quit")

	finalView := lipgloss.JoinVertical(lipgloss.Top, row1, footer)
	return containerStyle.Render(finalView)
}

func (ui *UI) SetContext(ctx *app.Context) {
	ui.ctx = ctx
}

func (ui *UI) SetCancel(cancel func()) {
	ui.cancel = cancel
}

func New() *UI {
	var ui UI
	ui.logs = NewLogView()
	ui.services = NewServiceView()

	return &ui
}
