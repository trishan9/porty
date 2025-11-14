package tui

import (
    "github.com/charmbracelet/lipgloss"
)

func GradientText(s string, cols []lipgloss.Color) string {
    if len(cols) == 0 || len(s) == 0 {
        return s
    }

    step := float64(len(cols)-1) / float64(len(s)-1)

    out := ""
    for i, r := range []rune(s) {
        ci := int(float64(i) * step)
        color := cols[ci]
        out += lipgloss.NewStyle().Foreground(color).Render(string(r))
    }

    return out
}
