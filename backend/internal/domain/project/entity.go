package project

import (
	"errors"
	"time"
)

var ErrEmptyName = errors.New("project name cannot be empty")

type Project struct {
	ID          string
	Name        string
	Description string
	CreatorID   string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewProject(id, name, creatorID string) (*Project, error) {
	if name == "" {
		return nil, ErrEmptyName
	}
	now := time.Now().UTC()
	return &Project{
		ID:        id,
		Name:      name,
		CreatorID: creatorID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
