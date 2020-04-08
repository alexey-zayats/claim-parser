package repository

import (
	"github.com/alexey-zayats/claim-parser/internal/database"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// RequestRepository ...
type RequestRepository struct {
	db *sqlx.DB
}

// RequestRepositoryInput ...
type RequestRepositoryInput struct {
	dig.In
	DB *sqlx.DB
}

// NewRequestRepository ...
func NewRequestRepository(param RequestRepositoryInput) interfaces.RequestRepository {
	return &RequestRepository{
		db: param.DB,
	}
}

// Create ...
func (r *RequestRepository) Create(data *model.Request) (int64, error) {

	var id int64

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		query := "INSERT INTO bids (" +
			"file_id, status, workflow_status, code, district, type, created_at, created_by, user_id, source" +
			") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

		res, err := t.Exec(query,
			data.FileID,
			data.Status,
			data.WorkflowStatus,
			data.Code,
			data.District,
			data.PassType,
			data.CreatedAt,
			data.CreatedBy,
			data.UserID,
			data.Source)

		if err != nil {
			return errors.Wrap(err, "unable create bid")
		}

		id, err = res.LastInsertId()
		if err != nil {
			return errors.Wrap(err, "unable get bid bids lastInsertID")
		}

		return nil
	})

	if err != nil {
		return 0, errors.Wrap(err, "transaction error")
	}

	return id, nil
}

// Update ...
func (r *RequestRepository) Update(data *model.Request) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "UPDATE bids SET file_id = ?, status = ?, workflow_status = ?, code = ?, " +
			"district = ?, type = ?, created_at = ?, created_by = ?, user_id = ?, source = ? " +
			"WHERE id = ?"
		_, err := t.Exec(sql, data.FileID,
			data.Status,
			data.WorkflowStatus,
			data.Code,
			data.District,
			data.PassType,
			data.CreatedAt,
			data.CreatedBy,
			data.UserID,
			data.ID,
			data.Source)

		if err != nil {
			return errors.Wrap(err, "unable update bids")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}

// Read ...
func (r *RequestRepository) Read(id int) (*model.Request, error) {
	var request *model.Request

	err := r.db.Get(request, "select * from bids where id=?", id)
	if err != nil {
		return nil, errors.Wrapf(err, "unable get bids record by id %s", id)
	}

	return request, nil
}

// Delete ...
func (r *RequestRepository) Delete(id int) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "DELETE FROM bids WHERE id = ?"
		_, err := t.Exec(sql, id)
		if err != nil {
			return errors.Wrapf(err, "unable delete from bids by id %d", id)
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}
