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
	textColor         = lipgloss.Color("#FAFAFA")
	selectedItemStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF5FAF"))

	mainStyle   = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 2, 0, 1)
	headerStyle = lipgloss.NewStyle().Bold(true).Foreground(textColor).Padding(0, 1)
	itemStyle   = lipgloss.NewStyle().Foreground(textColor)
)

type model struct {
	choices   []database.Task
	db        *database.Database
	cursor    int
	width     int
	height    int
	mode      Mode
	textInput textinput.Model
}

func initialModel(db *database.Database) *model {
	textInput := textinput.New()
	textInput.Placeholder = "Enter new todo..."
	textInput.CharLimit = 100
	textInput.Width = 30

	var tasks []database.Task
	tx := db.Find(&tasks)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}

	return &model{
		choices:   tasks,
		db:        db,
		textInput: textInput,
		mode:      modeList,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var tasks []database.Task
	tx := m.db.Find(&tasks)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}
	m.choices = tasks

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
		sort.Slice(m.choices, func(i, j int) bool {
			return m.choices[i].CreatedAt.Before(m.choices[j].CreatedAt)
		})

		header := headerStyle.Render("Tasks:")
		var items []string
		for i, choice := range m.choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">" // TODO: from config
			}

			checked := " "
			if choice.Done {
				checked = "x" // TODO: from config
			}

			itemText := fmt.Sprintf("%s [%s] %s", cursor, checked, choice.Name)
			var item string
			if m.cursor == i {
				item = selectedItemStyle.Render(itemText)
			} else {
				item = itemStyle.Render(itemText)
			}
			items = append(items, item)
		}

		// TODO: keybinds should be from config
		instructions := "Press `q` to quit | Press `a` to add a new todo | Press `d` to remove a todo"
		view := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			lipgloss.JoinVertical(lipgloss.Left, items...),
			instructions,
		)

		return m.float(view)
	case modeAdd:
		view := "Add New TODO:\n\n"
		view += m.textInput.View() + "\n\n"
		view += "Press Enter to confirm, Esc to cancel.\n"

		return m.float(view)
	}

	return "Unknown Mode"
}

func (m model) float(view string) string {
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, mainStyle.Render(view))
}

func (m model) mapAddModeActions(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		if input := m.textInput.Value(); input != "" {
			tx := m.db.DB.Save(&database.Task{Name: input})
			if tx.Error != nil {
				log.Println(tx.Error)
			}
		}

		m.mode = modeList
		return m, nil
	case "ctrl+c", "esc":
		m.mode = modeList
		return m, nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) toggleTaskMark() {
	if m.choices[m.cursor].Done {
		m.choices[m.cursor].Done = false
	} else {
		m.choices[m.cursor].Done = true
	}
	m.db.Save(&m.choices[m.cursor])
}

func (m *model) decementCursor() {
	if m.cursor > 0 {
		m.cursor--
	}
}

func (m *model) incrementCursor() {
	if m.cursor < len(m.choices)-1 {
		m.cursor++
	}
}

func (m *model) addTask() {
	m.mode = modeAdd
	m.textInput.Focus()
}

func (m *model) deleteTask() (tea.Model, tea.Cmd) {
	m.db.Delete(&m.choices[m.cursor])
	if m.cursor > 0 {
		m.cursor--
	}
	return m, func() tea.Msg { // NOTE: this is used to force a screen update
		return tea.WindowSizeMsg{Width: m.width, Height: m.height}
	}
}

func (m *model) updateWindowSize(msg tea.WindowSizeMsg) {
	m.width = msg.Width
	m.height = msg.Height
}
