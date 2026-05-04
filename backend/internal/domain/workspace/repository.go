package workspace

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Workspace, error)
	FindBySlug(ctx context.Context, slug string) (*Workspace, error)
	FindAll(ctx context.Context) ([]*Workspace, error)
	Save(ctx context.Context, ws *Workspace) error
	Update(ctx context.Context, ws *Workspace) error
	Delete(ctx context.Context, id string) error
}
