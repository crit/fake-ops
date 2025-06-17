package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
	"github.com/crit/fake-ops/internal/app"
)

type LogView struct {
	logs       []app.Message
	width      int
	height     int
	titleStyle lipgloss.Style
	infoStyle  lipgloss.Style
	errStyle   lipgloss.Style
	logStyle   lipgloss.Style
}

func NewLogView() *LogView {
	const numLastResults = 8

	return &LogView{
		logs:       make([]app.Message, numLastResults),
		titleStyle: lipgloss.NewStyle().Foreground(cPrimary).Bold(true),
		infoStyle:  lipgloss.NewStyle().Foreground(cSecondary),
		errStyle:   lipgloss.NewStyle().Foreground(cDanger),
		logStyle:   lipgloss.NewStyle().Foreground(cSecondary),
	}
}

func (lv *LogView) Update(msg app.Message) {
	lv.logs = append(lv.logs[1:], msg)
}

func (lv *LogView) Resize(width, height int) {
	lv.width = width
	lv.height = height
}

func (lv *LogView) View() string {
	var formatted []string

	formatted = append(formatted, lv.titleStyle.Render("LOGS"))
	for _, log := range lv.logs {
		switch log.Kind {
		case app.ErrorKind:
			formatted = append(formatted, lv.errStyle.Render(fmt.Sprintf("> %s", log)))
		case app.InfoKind:
			formatted = append(formatted, lv.infoStyle.Render(fmt.Sprintf("> %s", log)))
		}
	}

	return lv.logStyle.
		Height(lv.height).
		Width(lv.width).
		Render(lipgloss.JoinVertical(lipgloss.Top, formatted...))
}
