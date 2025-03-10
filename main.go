package main

import (
	"log"
	"os"

	"github.com/EwanGreer/todo-cli/config"
	"github.com/EwanGreer/todo-cli/database"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if os.Getenv("ENV") == "development" {
		f, err := tea.LogToFile("logs.log", "debug |")
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
	}

	cfg := config.NewConfig()

	db, err := database.NewDatabase(cfg.Database.Name)
	if err != nil {
		log.Fatal(err)
	}

	db.DB.Save(&database.List{
		Name: "Super Cool List",
		Tasks: []database.Task{
			{
				Name:        "A Task",
				Description: "Something to do.",
				Status:      database.InProgress,
			},
		},
	})

	p := tea.NewProgram(NewModel(db, cfg), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
