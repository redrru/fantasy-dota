package repository

import (
	"context"

	"github.com/redrru/fantasy-dota/internal/fantasy-dota/entity"
	"github.com/redrru/fantasy-dota/pkg/tracing"
)

func (r *Repository) ExampleList(ctx context.Context) ([]entity.ExampleModel, error) {
	ctx, span := tracing.DefaultTracer().Start(ctx, "ExampleList")
	defer span.End()

	var models []entity.ExampleModel

	err := r.db.Gorm.WithContext(ctx).Find(&models).Error
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (r *Repository) ExampleCreate(ctx context.Context, model entity.ExampleModel) error {
	ctx, span := tracing.DefaultTracer().Start(ctx, "ExampleCreate")
	defer span.End()

	return r.db.Gorm.WithContext(ctx).Create(&model).Error
}
