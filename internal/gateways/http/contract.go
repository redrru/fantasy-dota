package http

import (
	"context"

	"github.com/redrru/fantasy-dota/internal/fantasy-dota/entity"
)

type usecase interface {
	ExampleGet(ctx context.Context) ([]entity.ExampleModel, error)
	ExamplePost(ctx context.Context, model entity.ExampleModel) error
}
