package repository

import (
	"database/sql"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// SourceRepository ...
type SourceRepository struct {
	db *sqlx.DB
}

// SourceRepositoryInput ...
type SourceRepositoryInput struct {
	dig.In
	DB *sqlx.DB
}

// NewSourceRepository ...
func NewSourceRepository(param FilesRepositoryInput) interfaces.SourceRepository {
	return &SourceRepository{
		db: param.DB,
	}
}

// FindByName returns source record by name
func (r *SourceRepository) FindByName(name string) (*entity.Source, error) {
	var record entity.Source

	query := "SELECT id, name, title FROM sources WHERE name = ?"

	err := r.db.Get(&record, query, name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "unable find source record by name %s", name)
	}

	return &record, nil
}
