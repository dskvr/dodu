// Package tui implements the Bubbletea-based interactive disk-atlas UI.
package tui

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/tyutyutyu/dodu/pkg/docker"
	"github.com/tyutyutyu/dodu/pkg/group"
	"github.com/tyutyutyu/dodu/pkg/scan"
	"github.com/tyutyutyu/dodu/pkg/size"
)

// Layout is which root grouping the atlas shows.
type Layout int

// Available layouts.
const (
	LayoutByType Layout = iota
	LayoutByProject
)

// Model is the top-level Bubbletea model for dodu's TUI.
type Model struct {
	client  docker.Client
	logger  *slog.Logger
	scanner *scan.Scanner

	snap   *scan.Snapshot
	layout Layout
	sortBy group.SortKey
	root   *group.Node

	// Navigation: stack of "current node" frames. Top of stack is what we render.
	stack []*frame
	help  bool

	width, height int
	loading       bool
	loadStarted   time.Time
	err           error
	status        string
}

type frame struct {
	node     *group.Node
	selected int
	offset   int
}

type scanDoneMsg struct {
	snap *scan.Snapshot
	err  error
}

// NewModel constructs a new TUI model.
func NewModel(client docker.Client, logger *slog.Logger) *Model {
	if logger == nil {
		logger = slog.Default()
	}
	return &Model{
		client:  client,
		logger:  logger,
		scanner: scan.New(client, logger),
		layout:  LayoutByType,
		sortBy:  group.SortBySize,
		loading: true,
	}
}

// Init kicks off the initial scan.
func (m *Model) Init() tea.Cmd {
	m.loadStarted = time.Now()
	return m.scanCmd()
}

func (m *Model) scanCmd() tea.Cmd {
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		snap, err := m.scanner.Scan(ctx)
		return scanDoneMsg{snap: snap, err: err}
	}
}

// Update handles input and async messages.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case scanDoneMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
			return m, nil
		}
		m.snap = msg.snap
		m.rebuildTree()
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	}
	return m, nil
}

func (m *Model) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "?":
		m.help = !m.help
		return m, nil
	case "r":
		m.loading = true
		m.loadStarted = time.Now()
		m.err = nil
		m.status = "rescanning…"
		return m, m.scanCmd()
	case "g":
		if m.layout == LayoutByType {
			m.layout = LayoutByProject
		} else {
			m.layout = LayoutByType
		}
		m.rebuildTree()
		return m, nil
	case "s":
		switch m.sortBy {
		case group.SortBySize:
			m.sortBy = group.SortByName
		case group.SortByName:
			m.sortBy = group.SortByCount
		default:
			m.sortBy = group.SortBySize
		}
		if m.root != nil {
			m.root.Sort(m.sortBy)
		}
		return m, nil
	case "j", "down":
		m.move(1)
	case "k", "up":
		m.move(-1)
	case "enter", "l", "right":
		m.descend()
	case "h", "left", "esc", "backspace":
		m.ascend()
	case "g g", "home":
		f := m.top()
		if f != nil {
			f.selected = 0
			f.offset = 0
		}
	}
	return m, nil
}

func (m *Model) move(delta int) {
	f := m.top()
	if f == nil || f.node == nil {
		return
	}
	n := len(f.node.Children)
	if n == 0 {
		return
	}
	f.selected += delta
	if f.selected < 0 {
		f.selected = 0
	}
	if f.selected >= n {
		f.selected = n - 1
	}
}

func (m *Model) descend() {
	f := m.top()
	if f == nil || f.node == nil || len(f.node.Children) == 0 {
		return
	}
	child := f.node.Children[f.selected]
	if len(child.Children) == 0 {
		return
	}
	m.stack = append(m.stack, &frame{node: child})
}

func (m *Model) ascend() {
	if len(m.stack) > 1 {
		m.stack = m.stack[:len(m.stack)-1]
	}
}

func (m *Model) top() *frame {
	if len(m.stack) == 0 {
		return nil
	}
	return m.stack[len(m.stack)-1]
}

func (m *Model) rebuildTree() {
	if m.snap == nil {
		return
	}
	if m.layout == LayoutByType {
		m.root = group.ByType(m.snap)
	} else {
		m.root = group.ByProject(m.snap)
	}
	m.root.Sort(m.sortBy)
	m.stack = []*frame{{node: m.root}}
}

// View renders the current state.
func (m *Model) View() string {
	if m.loading {
		return loadingView(m.loadStarted)
	}
	if m.err != nil {
		return errStyle.Render(fmt.Sprintf("scan failed: %v\n\nq quit  r retry", m.err))
	}
	if m.help {
		return helpView()
	}
	return m.atlasView()
}

func (m *Model) atlasView() string {
	header := m.headerView()
	footer := m.footerView()
	bodyHeight := m.height - lipgloss.Height(header) - lipgloss.Height(footer)
	if bodyHeight < 5 {
		bodyHeight = 5
	}
	body := m.listView(bodyHeight)
	return lipgloss.JoinVertical(lipgloss.Left, header, body, footer)
}

func (m *Model) headerView() string {
	layout := "by-type"
	if m.layout == LayoutByProject {
		layout = "by-project"
	}
	sortName := []string{"size", "name", "count"}[m.sortBy]
	path := []string{}
	for _, f := range m.stack {
		if f.node != nil {
			path = append(path, f.node.Name)
		}
	}
	totals := ""
	if m.snap != nil && m.root != nil {
		totals = fmt.Sprintf("  total: %s", size.Format(m.root.Size.Total, size.IEC))
	}
	return headerStyle.Render(fmt.Sprintf("dodu  %s  sort:%s%s\n%s",
		layout, sortName, totals, strings.Join(path, " › ")))
}

func (m *Model) footerView() string {
	hint := "j/k move  l/enter open  h/esc back  g layout  s sort  r rescan  ? help  q quit"
	if m.status != "" {
		hint = m.status + "  •  " + hint
	}
	return footerStyle.Render(hint)
}

func (m *Model) listView(height int) string {
	f := m.top()
	if f == nil || f.node == nil || len(f.node.Children) == 0 {
		return bodyStyle.Render("(empty)")
	}
	children := f.node.Children
	// Adjust offset window to keep selected visible.
	if f.selected < f.offset {
		f.offset = f.selected
	}
	if f.selected >= f.offset+height {
		f.offset = f.selected - height + 1
	}
	if f.offset < 0 {
		f.offset = 0
	}
	end := f.offset + height
	if end > len(children) {
		end = len(children)
	}

	maxSize := int64(1)
	for _, c := range children {
		if c.Size.Exclusive > maxSize {
			maxSize = c.Size.Exclusive
		}
		if c.Size.Total > maxSize {
			maxSize = c.Size.Total
		}
	}

	rows := make([]string, 0, end-f.offset)
	for i := f.offset; i < end; i++ {
		c := children[i]
		row := renderRow(c, maxSize, i == f.selected)
		rows = append(rows, row)
	}
	return bodyStyle.Render(strings.Join(rows, "\n"))
}

func renderRow(n *group.Node, maxSize int64, selected bool) string {
	bar := bar20(n.Size.Total, maxSize)
	marker := " "
	if n.Size.Shared > 0 && n.Size.Total != n.Size.Exclusive {
		marker = "~"
	}
	if n.Size.Estimated {
		marker = "?"
	}
	line := fmt.Sprintf("%s %s %10s  %s",
		bar, marker, size.Format(n.Size.Total, size.IEC), n.Name)
	if selected {
		return selStyle.Render("▶ " + line)
	}
	return "  " + line
}

func bar20(v, max int64) string {
	const width = 20
	if max <= 0 {
		return strings.Repeat("░", width)
	}
	filled := int((v * int64(width)) / max)
	if filled < 0 {
		filled = 0
	}
	if filled > width {
		filled = width
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}

func loadingView(start time.Time) string {
	return headerStyle.Render(fmt.Sprintf("dodu — scanning Docker daemon… (%s)", time.Since(start).Truncate(time.Millisecond)))
}

func helpView() string {
	lines := []string{
		"dodu — keys",
		"",
		"  j / down      move down",
		"  k / up        move up",
		"  l / enter     open child",
		"  h / esc       back",
		"  g             toggle layout (by-type ↔ by-project)",
		"  s             cycle sort (size → name → count)",
		"  r             rescan",
		"  ?             toggle this help",
		"  q             quit",
	}
	return helpStyle.Render(strings.Join(lines, "\n"))
}
