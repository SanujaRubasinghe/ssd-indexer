package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/SanujaRubasinghe/ssdindexer/internal/scanner"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	stats         scanner.FileStats
	isDone        bool
	spinner       spinner.Model
	progress      progress.Model
	fileCount     int64
	totalFiles    int64
	ticker        *time.Ticker
}

type tickMsg time.Time
type UpdateMsg scanner.FileStats

func New() model {
	s := spinner.New()
	s.Spinner = spinner.Points

	p := progress.New(progress.WithDefaultGradient())
	return model{
		spinner:  s,
		progress: p,
		ticker:   time.NewTicker(120 * time.Millisecond),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		tickCmd(),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case spinner.TickMsg:
		if !m.isDone {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case tickMsg:
		if !m.isDone {
			return m, tickCmd()
		}

	case UpdateMsg:
		m.stats = scanner.FileStats(msg)
		// Update progress based on the scaled total (0-1000000)
		if msg.Total > 0 {
			percent := float64(msg.Total) / 1000000.0
			return m, m.progress.SetPercent(percent)
		}
		return m, nil

	case scanner.DoneMsg:
		m.isDone = true
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	title := lipgloss.NewStyle().Bold(true).
		Foreground(lipgloss.Color("#00E5FF")).
		Render("SSD Analyzer")

	if !m.isDone {
		// Calculate current progress percentage
		percent := float64(m.stats.Total) / 1000000.0 * 100
		return fmt.Sprintf(`
%s — Scanning...

%s  %.1f%% complete
`,
			title,
			m.spinner.View()+"  Scanning files...",
			percent,
		)
	}

	topExtensions := m.stats.GetTopOtherExtensions(5)
	extensionDetails := ""
	for _, ext := range topExtensions {
		if ext.Size > 0 {
			extPercent := pct(ext.Size, m.stats.Total) * 100
			extLabel := ext.Ext
			if extLabel == "" {
				extLabel = "(no extension)"
			}
			extensionDetails += fmt.Sprintf("    • %-12s %6.2f%%  %s\n", extLabel, extPercent, formatSize(ext.Size))
		}
	}

	// Get top 5 largest folders
	topFolders := m.stats.GetTopFolders(5)
	folderDetails := ""
	for i, folder := range topFolders {
		if folder.Size > 0 {
			// Get the last two components of the path for display
			displayPath := folder.Path
			if len(displayPath) > 40 {
				// If path is too long, show the last 40 characters with ellipsis
				displayPath = "..." + displayPath[max(0, len(displayPath)-40):]
			}
			folderDetails += fmt.Sprintf("    %d. %-40s  %s\n", 
				i+1, 
				displayPath, 
				formatSize(folder.Size))
		}
	}

	return fmt.Sprintf(`
%s — DONE

Memory Composition:

%s
%s
%s
%s
%s

Top 5 Largest Folders:
%s

Top Extensions (Other):
%s
Total Size: %.2f GB

Press Q to quit
`,
		title,
		m.bar("Photos", pct(m.stats.Photos, m.stats.Total), "#00FFAA", m.stats.Photos),
		m.bar("Videos", pct(m.stats.Videos, m.stats.Total), "#FF6A00", m.stats.Videos),
		m.bar("Docs  ", pct(m.stats.Docs, m.stats.Total), "#50A7FF", m.stats.Docs),
		m.bar("Compressed", pct(m.stats.Compressed, m.stats.Total), "#FF00FF", m.stats.Compressed),
		m.bar("Other ", pct(m.stats.Others, m.stats.Total), "#AAAAAA", m.stats.Others),
		folderDetails,
		extensionDetails,
		float64(m.stats.Total)/(1024*1024*1024),
	)
}

func pct(part, total int64) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total)
}

func (m model) bar(label string, percent float64, color string, size int64) string {
	barWidth := 30
	filled := int(percent * float64(barWidth))
	full := lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Render("=")

	empty := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#333333")).
		Render("-")

	// Format size in appropriate units
	sizeStr := formatSize(size)

	return fmt.Sprintf(
		"%s  %s%s  %.2f%%  %s",
		label,
		strings.Repeat(full, filled),
		strings.Repeat(empty, barWidth-filled),
		percent*100,
		sizeStr,
	)
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%7.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func tickCmd() tea.Cmd {
	return func() tea.Msg { return tickMsg(time.Now()) }
}
