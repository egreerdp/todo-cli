package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/EwanGreer/todo-cli/database"
	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Mode int

const (
	modeList Mode = iota
	modeAdd
)

var (
// textColor         = lipgloss.Color("#FAFAFA")
// selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5FAF"))
//
// mainStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 2, 0, 1)
// headerStyle = lipgloss.NewStyle().Bold(true).Foreground(textColor).Padding(0, 1)
// itemStyle   = lipgloss.NewStyle().Foreground(textColor)
)

type model struct {
	Lists  []database.List
	db     *database.Database
	cursor int
	width  int
	height int
	mode   Mode
}

func initialModel(db *database.Database) *model {
	var lists []database.List
	if err := db.Preload("Tasks").Find(&lists).Error; err != nil {
		log.Fatal("Error fetching lists with tasks:", err)
	}

	return &model{
		Lists: lists,
		db:    db,
		mode:  modeList,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var lists []database.List
	if err := m.db.Preload("Tasks").Find(&lists).Error; err != nil {
		log.Fatal("Error fetching lists with tasks:", err)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case modeAdd:
			return m.mapAddModeActions(msg)
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.decementCursor()
		case "down", "j":
			m.incrementCursor()
		case "enter", " ", "x":
			m.toggleTaskMark()
		case "a":
			m.addTask()
		case "e":
			// edit task
		case "d":
			return m.deleteTask()
		case "?":
			// show help menu
		case ".":
			// show/hide completed tasks
		case "@":
			// show/hide command history
		}
	case tea.WindowSizeMsg:
		m.updateWindowSize(msg)
	}

	return m, nil
}

func (m model) View() string {
	switch m.mode {
	case modeList:
		sort.Slice(m.Lists, func(i, j int) bool {
			return m.Lists[i].CreatedAt.Before(m.Lists[j].CreatedAt)
		})

		var ListNames []string
		for i, list := range m.Lists {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			itemText := fmt.Sprintf("%s %s", cursor, list.Name)
			ListNames = append(ListNames, itemText)
		}

		// NOTE: I feel like I need a curor per scollable element on screen?
		// If this was the case, I would just need to track which element was focused, and use its cursor.
		// This might be dumb, but a cursor that can move up down left and right feels hard.
		var taskNames []string
		for i, task := range m.Lists[m.cursor].Tasks {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			itemText := fmt.Sprintf("%s %s", cursor, task.Name)
			taskNames = append(taskNames, itemText)
		}

		halfWidth := m.width / 4
		leftBox := lipgloss.NewStyle().
			Width(halfWidth).
			Height(m.height-3).
			MarginLeft(1).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			Align(lipgloss.Left, lipgloss.Top).
			Render(lipgloss.JoinVertical(lipgloss.Left, ListNames...))

		centerBox := lipgloss.NewStyle().
			Width(halfWidth).
			Height(m.height-3).
			MarginLeft(1).
			Padding(0, 1).
			Border(lipgloss.RoundedBorder()).
			Align(lipgloss.Left, lipgloss.Top).
			Render(lipgloss.JoinVertical(lipgloss.Left, taskNames...))

		return lipgloss.JoinHorizontal(lipgloss.Left, leftBox, centerBox)
	}
	return "Unknown Mode"
}

func (m model) mapAddModeActions(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		m.mode = modeList
		return m, nil
	case "ctrl+c", "esc":
		m.mode = modeList
		return m, nil
	}

	var cmd tea.Cmd
	return m, cmd
}

func (m model) toggleTaskMark() {
}

func (m *model) decementCursor() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *model) incrementCursor() {
	// if m.cursor < len(m.choices)-1 { // TODO: this will need to consider the focused block/element
	// 	m.cursor++
	// }
}

func (m *model) addTask() {
	m.mode = modeAdd
}

func (m *model) deleteTask() (tea.Model, tea.Cmd) {
	// m.db.Delete(&m.choices[m.cursor])
	// if m.cursor > 0 {
	// 	m.cursor--
	// }
	// return m, func() tea.Msg { // NOTE: this is used to force a screen update
	// 	return tea.WindowSizeMsg{Width: m.width, Height: m.height}
	// }
	return nil, nil
}

func (m *model) updateWindowSize(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.height = msg.Height
}
