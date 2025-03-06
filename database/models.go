package database

import (
	"gorm.io/gorm"
)

type Status int

const (
	Undefined Status = iota
	Todo
	InProgress
	Done
)

func (s Status) String() string {
	switch s {
	case 1:
		return "ToDo"
	case 2:
		return "In Progress"
	case 3:
		return "Done"
	}

	return "Undefined"
}

type List struct {
	gorm.Model
	Name  string
	Tasks []Task `gorm:"foreignKey:ListId"`
}

type Task struct {
	gorm.Model
	Name        string
	Description string
	Status      Status
	ListId      uint
	ParentId    *uint  `gorm:"index"`
	Subtasks    []Task `gorm:"foreignKey:ParentId"`
}
