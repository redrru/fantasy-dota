package usecase

import (
	"context"

	"github.com/redrru/fantasy-dota/internal/fantasy-dota/entity"
)

type repository interface {
	ExampleList(ctx context.Context) ([]entity.ExampleModel, error)
	ExampleCreate(ctx context.Context, model entity.ExampleModel) error
}
