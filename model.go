package main

import (
	"fmt"
	"log"

	"github.com/EwanGreer/todo-cli/components/list"
	"github.com/EwanGreer/todo-cli/config"
	"github.com/EwanGreer/todo-cli/database"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mainModel struct {
	Elements      []tea.Model
	focusedIndex  uint
	db            *database.Database
	Width, Height int
}

func NewModel(db *database.Database, cfg *config.Config) *mainModel {
	lists := []database.List{}
	tx := db.Find(&lists)
	if tx.Error != nil {
		log.Fatal(tx.Error)
	}

	return &mainModel{
		Elements: []tea.Model{
			&list.ListModel{
				Db:    db,
				Items: lists,
			},
		},
		focusedIndex: 0,
		db:           db,
	}
}

func (m *mainModel) Init() tea.Cmd { return nil }

func (m *mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab", "l":
			if m.focusedIndex < uint(len(m.Elements)-1) {
				m.focusedIndex++
			} else {
				m.focusedIndex = 0
			}
		case "shift+tab", "h":
			if m.focusedIndex > 0 {
				m.focusedIndex--
			} else {
				m.focusedIndex = uint(len(m.Elements) - 1)
			}
		case "a":
			var cmd tea.Cmd
			m.Elements[m.focusedIndex], cmd = m.Elements[m.focusedIndex].Update(msg) // returns a pointer, so why does it need reassigned?
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
	}

	return m, nil
}

func (m *mainModel) View() string {
	focusedStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("162"))
		// TODO: how to only have 1 border
	unfocusedStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("8"))

	views := []string{}

	focusedCounter := fmt.Sprintf("%d", m.focusedIndex)
	views = append(views, focusedCounter)

	for i, e := range m.Elements {
		if m.focusedIndex == uint(i) {
			views = append(views, focusedStyle.Render(e.View()))
			continue
		}

		views = append(views, unfocusedStyle.Render(e.View()))
	}

	return lipgloss.JoinHorizontal(lipgloss.Left, views...)
}
