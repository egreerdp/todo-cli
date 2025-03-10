package list

import (
	"github.com/EwanGreer/todo-cli/database"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ListModel struct {
	Name  string
	Items []database.List // this needs to be updated for bubbletea to rerender
	Db    *database.Database
}

func (l *ListModel) Init() tea.Cmd { return nil }

func (l *ListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			newList := database.List{
				Name: "Even Cooler List",
				Tasks: []database.Task{
					{
						Name:        "A Task",
						Description: "Something to do.",
						Status:      database.InProgress,
					},
				},
			}

			l.Db.DB.Save(&newList)

			var items []database.List
			_ = l.Db.Find(&items)
			l.Items = items

			return nil, nil
		}
	}
	return nil, nil
}

func (l *ListModel) View() string {
	style := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), true).Align(lipgloss.Center)

	itemNames := []string{}
	for _, item := range l.Items {
		itemNames = append(itemNames, item.Name)
	}

	return style.Render(lipgloss.JoinVertical(lipgloss.Left, itemNames...))
}
