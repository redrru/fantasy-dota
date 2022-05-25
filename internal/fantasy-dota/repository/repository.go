package repository

import (
	postgres "github.com/redrru/fantasy-dota/pkg/db"
)

type Repository struct {
	db *postgres.DB
}

func NewRepository(db *postgres.DB) *Repository {
	return &Repository{db: db}
}
