package tui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/trishan9/porty/internal"
)

// ---------- Tokyo Night palette ----------

var (
    // No backgrounds. Terminal decides the background.
    textColor    = lipgloss.Color("#c0caf5")
    mutedColor   = lipgloss.Color("#6b7089")

    successColor = lipgloss.Color("#9ece6a")
    warnColor    = lipgloss.Color("#e0af68")
    errorColor   = lipgloss.Color("#f7768e")

    blueColor    = lipgloss.Color("#7aa2f7")
    cyanColor    = lipgloss.Color("#7dcfff")
    purpleColor  = lipgloss.Color("#bb9af7")

    // gradients still fine
    gradientColors = []lipgloss.Color{
        cyanColor,
        blueColor,
        purpleColor,
    }
)

// Cursor highlight background (auto adjusts to terminal theme)
var cursorBg = lipgloss.AdaptiveColor{
    Light: "#d0d7e5",
    Dark:  "#2f3348",
}

var cursorFg = lipgloss.AdaptiveColor{
    Light: "#1a1b26",
    Dark:  "#c0caf5",
}


var (
	baseStyle = lipgloss.NewStyle().
		Foreground(textColor)

	panelStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(mutedColor).
		Padding(1, 2).
		Margin(1, 1)

	helpStyle = lipgloss.NewStyle().
		Foreground(mutedColor).
		Margin(0, 2)

	titleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(blueColor).
		Margin(0, 2)

	statusSuccess = lipgloss.NewStyle().Foreground(successColor).Margin(0, 2)
	statusError   = lipgloss.NewStyle().Foreground(errorColor).Margin(0, 2)
	statusNeutral = lipgloss.NewStyle().Foreground(mutedColor).Margin(0, 2)
)

const tickInterval = 2 * time.Second

type tickMsg struct{}

type cpuSample struct {
	idle  uint64
	total uint64
}

// ---------- model ----------

type model struct {
	entries  []internal.PortEntry
	cursor   int
	selected map[int]bool
	status   string
	statusOK bool

	cpuPercent  int
	memUsedMiB  int
	memTotalMiB int
	lastCPU     cpuSample
}

// NewModel creates the initial TUI model.
func NewModel(entries []internal.PortEntry) model {
	m := model{
		entries:  entries,
		cursor:   0,
		selected: make(map[int]bool),
		status:   "↑/↓/j/k move  space select  enter/x kill  r reload  q quit",
		statusOK: true,
	}
	m = refreshModel(m) // initial stats/ports snapshot
	return m
}

// Run launches the Bubble Tea program.
func Run(entries []internal.PortEntry) error {
	p := tea.NewProgram(NewModel(entries))
	_, err := p.Run()
	return err
}

func (m model) Init() tea.Cmd { return tickCmd() }

func tickCmd() tea.Cmd {
	return tea.Tick(tickInterval, func(time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tickMsg:
		m = refreshModel(m)
		return m, tickCmd()

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "esc", "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.entries)-1 {
				m.cursor++
			}

		case " ":
			if len(m.entries) == 0 {
				return m, nil
			}
			m.selected[m.cursor] = !m.selected[m.cursor]

		case "r":
			m = refreshModel(m)
			m.status = "reloaded"
			m.statusOK = true

		case "enter", "x":
			if len(m.entries) == 0 {
				return m, nil
			}
			pids := m.collectSelectedPIDs()
			if len(pids) == 0 {
				if pid := m.entries[m.cursor].PID; pid > 0 {
					pids = []int{pid}
				}
			}
			if len(pids) == 0 {
				m.status = "no valid PIDs to kill"
				m.statusOK = false
				return m, nil
			}

			msgs := internal.KillPIDs(pids)
			m.status = strings.Join(msgs, " | ")

			m.statusOK = true
			for _, s := range msgs {
				ls := strings.ToLower(s)
				if strings.Contains(ls, "fail") || strings.Contains(ls, "error") {
					m.statusOK = false
					break
				}
			}
		}
	}
	return m, nil
}

func (m model) collectSelectedPIDs() []int {
	var pids []int
	for idx, sel := range m.selected {
		if sel && idx >= 0 && idx < len(m.entries) {
			if pid := m.entries[idx].PID; pid > 0 {
				pids = append(pids, pid)
			}
		}
	}
	return pids
}

// ---------- refresh ----------

func refreshModel(m model) model {
	// refresh ports
	if entries, err := internal.ListPorts(); err == nil {
		m.entries = entries
	}

	// refresh memory
	used, total := readMem()
	if total > 0 {
		m.memUsedMiB = used
		m.memTotalMiB = total
	}

	// refresh cpu
	pct, sample := readCPU(m.lastCPU)
	if sample.total != 0 {
		m.cpuPercent = pct
		m.lastCPU = sample
	}

	return m
}

// ---------- system info ----------

func readMem() (usedMiB, totalMiB int) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return
	}
	lines := strings.Split(string(data), "\n")
	var memTotal, memAvailable uint64

	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			f := strings.Fields(line)
			if len(f) >= 2 {
				memTotal, _ = strconv.ParseUint(f[1], 10, 64)
			}
		}
		if strings.HasPrefix(line, "MemAvailable:") {
			f := strings.Fields(line)
			if len(f) >= 2 {
				memAvailable, _ = strconv.ParseUint(f[1], 10, 64)
			}
		}
	}

	if memTotal == 0 {
		return
	}

	totalMiB = int(memTotal / 1024)
	usedMiB = int((memTotal - memAvailable) / 1024)
	return
}

func readCPU(prev cpuSample) (int, cpuSample) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, prev
	}
	line := strings.SplitN(string(data), "\n", 2)[0] // "cpu ..."
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return 0, prev
	}

	var total uint64
	for i := 1; i < len(fields); i++ {
		v, _ := strconv.ParseUint(fields[i], 10, 64)
		total += v
	}
	idle, _ := strconv.ParseUint(fields[4], 10, 64)
	cur := cpuSample{idle: idle, total: total}

	if prev.total == 0 {
		return 0, cur
	}

	deltaTotal := float64(cur.total - prev.total)
	deltaIdle := float64(cur.idle - prev.idle)
	if deltaTotal <= 0 {
		return 0, cur
	}
	usage := int((1.0 - deltaIdle/deltaTotal) * 100.0)
	if usage < 0 {
		usage = 0
	}
	if usage > 100 {
		usage = 100
	}
	return usage, cur
}

// ---------- view ----------

func (m model) View() string {
	portsPanel := m.renderPortsPanel()

	// vertical layout works nicely on most widths; lipgloss handles wrapping
	main := lipgloss.JoinVertical(lipgloss.Left, portsPanel)

	var statusLine string
	if m.status == "" {
		statusLine = statusNeutral.Render("ready")
	} else if m.statusOK {
		statusLine = statusSuccess.Render(m.status)
	} else {
		statusLine = statusError.Render(m.status)
	}

	help := helpStyle.Render("↑/↓/j/k move  space select  enter/x kill  r reload  q quit")

	return baseStyle.Render(
		titleStyle.Render("PORTY – Listening Ports") + "\n\n" +
			main + "\n" + help + "\n" + statusLine + "\n",
	)
}

func (m model) renderPortsPanel() string {
	if len(m.entries) == 0 {
		return panelStyle.Render("No listening ports detected.")
	}

	var b strings.Builder

	header := gradientText(" LISTENING PORTS ", gradientColors)
	b.WriteString(header + "\n\n")

	// table header
	headerLine := fmt.Sprintf("  %-3s %-2s %-7s %-6s %-6s %-22s %-8s %-12s %-8s",
		"#", " ", "STATE", "PORT", "PROTO", "PROCESS", "PID", "USER", "TAG")
	b.WriteString(headerLine + "\n")
	b.WriteString(strings.Repeat("─", len(headerLine)) + "\n")

	for i, e := range m.entries {
		cursor := " "
		if i == m.cursor {
			cursor = "▸"
		}
		check := "○"
		if m.selected[i] {
			check = "●"
		}

		idxStr := fmt.Sprintf("%2d", i+1)

		portStr := gradientText(fmt.Sprintf("%-6s", e.LocalPort), gradientColors)
		protoStr := lipgloss.NewStyle().Foreground(blueColor).Render(e.Proto)
		stateStr := lipgloss.NewStyle().Foreground(mutedColor).Render(e.State)

		proc := truncate(e.ProcessName, 22)
		user := truncate(e.UserName, 12)

		tagText, tagStyle := styleTag(e.Tag)
		tagRendered := tagStyle.Render(tagText)

		pidStr := "-"
		if e.PID > 0 {
			pidStr = fmt.Sprintf("%d", e.PID)
		}

		row := fmt.Sprintf("  %-3s %s %-7s %s %-6s %-22s %-8s %-12s %-8s",
			idxStr, check, stateStr, portStr, protoStr, proc, pidStr, user, tagRendered)

		if i == m.cursor {
			row = lipgloss.NewStyle().
				Background(cursorBg).
				Foreground(cursorFg).
				Render(cursor + " " + row)
		} else {
			row = "  " + cursor + " " + row
		}


		b.WriteString(row + "\n")
	}

	return panelStyle.Render(b.String())
}

// ---------- helpers ----------

func bar(percent, width int) string {
	if percent < 0 {
		percent = 0
	}
	if percent > 100 {
		percent = 100
	}
	filled := int(float64(percent) / 100.0 * float64(width))
	if filled > width {
		filled = width
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

func styleTag(tag string) (string, lipgloss.Style) {
    switch strings.ToUpper(tag) {
    case "USER":
        return "USER", lipgloss.NewStyle().Foreground(successColor)
    case "SYSTEM":
        return "SYSTEM", lipgloss.NewStyle().Foreground(warnColor)
    case "SELF":
        return "SELF", lipgloss.NewStyle().Foreground(cyanColor)
    case "KERNEL":
        return "KERNEL", lipgloss.NewStyle().Foreground(purpleColor).Bold(true)
    default:
        return "UNKNOWN", lipgloss.NewStyle().Foreground(errorColor)
    }
}

func gradientText(s string, cols []lipgloss.Color) string {
	runes := []rune(s)
	if len(runes) == 0 || len(cols) == 0 {
		return s
	}
	if len(runes) == 1 {
		return lipgloss.NewStyle().Foreground(cols[0]).Render(s)
	}

	step := float64(len(cols)-1) / float64(len(runes)-1)
	var out strings.Builder

	for i, r := range runes {
		ci := int(float64(i) * step)
		if ci < 0 {
			ci = 0
		}
		if ci >= len(cols) {
			ci = len(cols) - 1
		}
		out.WriteString(
			lipgloss.NewStyle().Foreground(cols[ci]).Render(string(r)),
		)
	}
	return out.String()
}

func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	if max <= 3 {
		return string(r[:max])
	}
	return string(r[:max-3]) + "..."
}
