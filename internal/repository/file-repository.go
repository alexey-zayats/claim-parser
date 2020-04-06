package repository

import (
	"github.com/alexey-zayats/claim-parser/internal/database"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
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

// Update ...
func (r *FileRepository) Update(data *model.File) error {

	logrus.WithFields(logrus.Fields{"data": data}).Debug("FileRepository.Update")

	err := database.WithTransaction(r.db, func(t database.Transaction) error {
		_, err := t.Exec("UPDATE files SET status = ?, log = ? WHERE id = ?", data.Status, data.Log, data.ID)
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
