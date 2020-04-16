package repository

import (
	"github.com/alexey-zayats/claim-parser/internal/database"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// BranchRepository ...
type BranchRepository struct {
	db *sqlx.DB
}

// BranchRepositoryInput ...
type BranchRepositoryInput struct {
	dig.In
	DB *sqlx.DB
}

// NewBranchRepository ...
func NewBranchRepository(param FilesRepositoryInput) interfaces.BranchRepository {
	return &BranchRepository{
		db: param.DB,
	}
}

// Create ...
func (r *BranchRepository) Create(data *entity.Branch) error {
	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		query := "INSERT INTO branches (name, type) VALUES (?, ?)"
		res, err := t.Exec(query, data.Name, data.Type)
		if err != nil {
			return err
		}

		data.ID, err = res.LastInsertId()
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}

// GetAll ...
func (r *BranchRepository) GetAll() ([]*entity.Branch, error) {

	list := make([]*entity.Branch, 0)

	query := "SELECT id, name, type FROM branches"
	rows, err := r.db.Queryx(query)
	if err != nil {
		return nil, errors.Wrapf(err, "unable get branches")
	}

	for rows.Next() {
		var branch entity.Branch
		if err := rows.StructScan(&branch); err != nil {
			return nil, errors.Wrapf(err, "unable scan row")
		}
		list = append(list, &branch)
	}

	return list, nil
}
