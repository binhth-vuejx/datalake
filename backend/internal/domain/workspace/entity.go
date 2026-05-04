package workspace

import (
	"errors"
	"time"
)

var ErrEmptySlug = errors.New("workspace slug cannot be empty")

type Workspace struct {
	ID          string
	Name        string
	Slug        string
	Description string
	OwnerID     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewWorkspace(id, name, slug, ownerID string) (*Workspace, error) {
	if slug == "" {
		return nil, ErrEmptySlug
	}
	now := time.Now().UTC()
	return &Workspace{
		ID:        id,
		Name:      name,
		Slug:      slug,
		OwnerID:   ownerID,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}
