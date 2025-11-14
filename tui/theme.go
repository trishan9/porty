package tui

import "github.com/charmbracelet/lipgloss"

// Tokyo Night Palette
var (
    TNBackground = lipgloss.Color("#1a1b26")
    TNPanel      = lipgloss.Color("#24283b")

    TNText       = lipgloss.Color("#c0caf5")
    TNSuccess    = lipgloss.Color("#9ece6a")
    TNWarn       = lipgloss.Color("#e0af68")
    TNError      = lipgloss.Color("#f7768e")

    TNBlue       = lipgloss.Color("#7aa2f7")
    TNPurple     = lipgloss.Color("#bb9af7")
    TNCyan       = lipgloss.Color("#7dcfff")

    TNGradient   = []lipgloss.Color{
        lipgloss.Color("#7dcfff"),
        lipgloss.Color("#7aa2f7"),
        lipgloss.Color("#bb9af7"),
    }
)

var (
    TitleStyle = lipgloss.NewStyle().
        Foreground(TNBlue).
        Bold(true).
        Padding(0, 1)

    SelectedStyle = lipgloss.NewStyle().
        Background(TNPanel).
        Foreground(TNText)

    BorderStyle = lipgloss.NewStyle().
        Border(lipgloss.RoundedBorder()).
        BorderForeground(TNPanel).
        Padding(1, 2)

    FadedStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("#565f89"))
)
