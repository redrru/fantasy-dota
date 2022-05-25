package usecase

import (
	"context"
	"fmt"

	"github.com/redrru/fantasy-dota/internal/fantasy-dota/entity"
)

func (u *Usecase) ExampleGet(ctx context.Context) ([]entity.ExampleModel, error) {
	models, err := u.repo.ExampleList(ctx)
	if err != nil {
		return nil, err
	}

	if models == nil {
		return nil, fmt.Errorf("empty models")
	}

	return models, err
}

func (u *Usecase) ExamplePost(ctx context.Context, model entity.ExampleModel) error {
	return u.repo.ExampleCreate(ctx, model)
}
