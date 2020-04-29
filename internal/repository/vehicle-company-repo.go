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

// VehicleCompanyRepo ...
type VehicleCompanyRepo struct {
	db *sqlx.DB
}

// VehicleCompanyRepoDI ...
type VehicleCompanyRepoDI struct {
	dig.In
	DB *sqlx.DB
}

// NewVehicleCompanyRepo ...
func NewVehicleCompanyRepo(di VehicleCompanyRepoDI) interfaces.VehicleCompanyRepo {
	return &VehicleCompanyRepo{
		db: di.DB,
	}
}

// FindByOgrnInn ...
func (r *VehicleCompanyRepo) FindByOgrnInn(ogrn, inn string) (*entity.Company, error) {
	var record entity.Company

	query :=
		"SELECT id, ogrn, inn, name, branch_id, status " +
			"FROM companies " +
			"WHERE ogrn = ? AND inn = ?"

	err := r.db.Get(&record, query, ogrn, inn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "unable get company_people record by inn(%d) and ogrn(%d)", inn)
	}

	return &record, nil
}

// FindByINN ...
func (r *VehicleCompanyRepo) FindByINN(inn string) (*entity.Company, error) {
	var record entity.Company

	query := "SELECT id, ogrn, inn, name, branch_id, status " +
		"FROM companies WHERE inn = ?"

	err := r.db.Get(&record, query, inn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "unable get company_people record by inn %s", inn)
	}

	return &record, nil
}

// FindByOGRN ...
func (r *VehicleCompanyRepo) FindByOGRN(ogrn string) (*entity.Company, error) {
	var record entity.Company

	query := "SELECT " +
		"id, ogrn, inn, name, branch_id, status " +
		"FROM companies " +
		"WHERE ogrn = ?"

	err := r.db.Get(&record, query, ogrn)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "unable get company_people record by ogrn %s", ogrn)
	}

	return &record, nil
}

// Create ...
func (r *VehicleCompanyRepo) Create(data *entity.Company) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		query := "INSERT INTO companies (" +
			"ogrn, inn, name, branch_id, status" +
			") VALUES (?, ?, ?, ?, ?)"

		res, err := t.Exec(query,
			data.OGRN,
			data.INN,
			data.Name,
			data.BranchID,
			data.Status,
		)

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

// Update ...
func (r *VehicleCompanyRepo) Update(data *entity.Company) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "UPDATE companies " +
			"SET " +
			"ogrn = ?, inn = ?, name = ?, branch_id = ?, status = ? " +
			"WHERE id = ?"

		_, err := t.Exec(sql,
			data.OGRN,
			data.INN,
			data.Name,
			data.BranchID,
			data.Status,
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
func (r *VehicleCompanyRepo) Read(id int64) (*entity.Company, error) {
	var pass entity.Company

	query :=
		"SELECT " +
			"id, " +
			"ogrn, inn, name, branch_id, status " +
			"FROM " +
			"companies " +
			"WHERE " +
			"id = ?"

	err := r.db.Get(&pass, query, id)
	if err != nil {
		return nil, errors.Wrapf(err, "unable get companies_people record by id %s", id)
	}

	return &pass, nil
}

// Delete ...
func (r *VehicleCompanyRepo) Delete(id int64) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "DELETE FROM companies WHERE id = ?"
		_, err := t.Exec(sql, id)
		if err != nil {
			return errors.Wrapf(err, "unable delete from companies_people by id %d", id)
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}
