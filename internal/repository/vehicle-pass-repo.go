package repository

import (
	"database/sql"
	"github.com/alexey-zayats/claim-parser/internal/database"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// VehiclePassRepo ...
type VehiclePassRepo struct {
	db *sqlx.DB
}

// VehiclePassRepoInput ...
type VehiclePassRepoInput struct {
	dig.In
	DB *sqlx.DB
}

// NewVehiclePassRepo ...
func NewVehiclePassRepo(param VehiclePassRepoInput) interfaces.VehiclePassRepo {
	return &VehiclePassRepo{
		db: param.DB,
	}
}

// FindByCar ...
func (r *VehiclePassRepo) FindByCar(car string) (*entity.Pass, error) {
	var record entity.Pass

	query := "SELECT id, bid_id, " +
		"lastname, firstname, patrname, car, " +
		"source, district_id, pass_type, pass_number, " +
		"shipping, status, file_id, " +
		"created_at, created_by " +
		"FROM passes WHERE car = ?"

	err := r.db.Get(&record, query, car)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "unable get passes record by car %s", car)
	}

	return &record, nil
}

// Create ...
func (r *VehiclePassRepo) Create(data *entity.Pass) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		query := "INSERT INTO passes (" +
			"bid_id, " +
			"lastname, firstname, patrname, car, " +
			"source, district_id, pass_type, pass_number, " +
			"shipping, status, file_id, " +
			"created_at, created_by" +
			") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

		res, err := t.Exec(query,
			data.BidID,
			data.Lastname,
			data.Firstname,
			data.Patrname,
			data.Car,
			data.Source,
			data.DistrictID,
			data.PassType,
			data.PassNumber,
			data.Shipping,
			data.Status,
			data.FileID,
			data.CreatedAt,
			data.CreatedBy,
		)

		if err != nil {
			return err
		}

		data.ID, err = res.LastInsertId()
		if err != nil {
			return errors.Wrap(err, "unable get passes lasInsertId")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}

// Update ...
func (r *VehiclePassRepo) Update(data *entity.Pass) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "UPDATE passes SET " +
			"bid_id = ?, " +
			"lastname = ?, firstname = ?, patrname = ?, car = ?, " +
			"source = ?, district_id = ?, pass_type = ?, pass_number = ?, " +
			"shipping = ?, status = ?, file_id = ?, " +
			"created_at = ?, created_by = ? " +
			"WHERE id = ?"

		_, err := t.Exec(sql,
			data.BidID,
			data.Lastname,
			data.Firstname,
			data.Patrname,
			data.Car,
			data.Source,
			data.DistrictID,
			data.PassType,
			data.PassNumber,
			data.Shipping,
			data.Status,
			data.FileID,
			data.CreatedAt,
			data.CreatedBy,
			data.ID)

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

// Read ...
func (r *VehiclePassRepo) Read(id int64) (*entity.Pass, error) {
	var pass entity.Pass

	query := "SELECT " +
		"id, bid_id, issued_id, " +
		"lastname, firstname, patrname, car, " +
		"source, district_id, pass_type, pass_number, " +
		"shipping, status, file_id, " +
		"created_at, created_by " +
		"FROM passes WHERE id = ?"

	err := r.db.Get(&pass, query, id)
	if err != nil {
		return nil, errors.Wrapf(err, "unable get passes record by id %s", id)
	}

	return &pass, nil
}

// Delete ...
func (r *VehiclePassRepo) Delete(id int64) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "DELETE FROM passes WHERE id = ?"
		_, err := t.Exec(sql, id)
		if err != nil {
			return errors.Wrapf(err, "unable delete from passes by id %d", id)
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}
