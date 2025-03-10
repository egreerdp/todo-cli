package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
}

func NewDatabase(name string) (*Database, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(homeDir, fmt.Sprintf("%s.db", name))
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&List{}, &Task{}); err != nil {
		return nil, err
	}

	return &Database{db}, nil
}

func (d *Database) Save(task *Task) {
	tx := d.DB.Save(task)
	if tx.Error != nil {
		log.Println(tx.Error)
	}
}
