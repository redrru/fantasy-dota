package repository

import (
	"context"

	"github.com/redrru/fantasy-dota/internal/fantasy-dota/entity"
)

func (r *Repository) ExampleList(ctx context.Context) ([]entity.ExampleModel, error) {
	var models []entity.ExampleModel

	err := r.db.Gorm.WithContext(ctx).Find(&models).Error
	if err != nil {
		return nil, err
	}

	return models, nil
}

func (r *Repository) ExampleCreate(ctx context.Context, model entity.ExampleModel) error {
	return r.db.Gorm.WithContext(ctx).Create(&model).Error
}
