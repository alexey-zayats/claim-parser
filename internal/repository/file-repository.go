package repository

import (
	"github.com/alexey-zayats/claim-parser/internal/database"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// FileRepository ...
type FileRepository struct {
	db *sqlx.DB
}

// FilesRepositoryInput ...
type FilesRepositoryInput struct {
	dig.In
	DB *sqlx.DB
}

// NewFileRepository ...
func NewFileRepository(param FilesRepositoryInput) interfaces.FileRepository {
	return &FileRepository{
		db: param.DB,
	}
}

// Create ...
func (r *FileRepository) Create(data *entity.File) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		query := "INSERT INTO files (filepath, status, log, created_at, source) VALUES (?, ?, ?, ?, ?)"

		res, err := t.Exec(query, data.Filepath, data.Status, data.Log, data.CreatedAt, data.Source)
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

// UpdateState ...
func (r *FileRepository) UpdateState(data *entity.File) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		query := "UPDATE files " +
			"SET " +
			"status = ?, log = concat(log, ?), source = ? " +
			"WHERE id = ?"

		_, err := t.Exec(query,
			data.Status,
			data.Log,
			data.Source,
			data.ID)

		if err != nil {
			return errors.Wrap(err, "unable update files")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}

// Read ...
func (r *FileRepository) Read(id int64) (*entity.File, error) {
	var file entity.File

	err := r.db.Get(&file, "SELECT "+
		"id, filepath, status, log, created_at, source "+
		"FROM files WHERE id = ?", id)
	if err != nil {
		return nil, errors.Wrapf(err, "unable get file record by id %s", id)
	}

	return &file, nil
}

// Delete ...
func (r *FileRepository) Delete(id int64) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "DELETE FROM files WHERE id = ?"
		_, err := t.Exec(sql, id)
		if err != nil {
			return errors.Wrapf(err, "unable delete from files by id %d", id)
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}
