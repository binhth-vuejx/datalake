package project

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (*Project, error)
	FindAll(ctx context.Context, limit, offset int) ([]*Project, int, error)
	Save(ctx context.Context, project *Project) error
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id string) error
}
