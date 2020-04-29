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

// VehicleIssuedRepo ...
type VehicleIssuedRepo struct {
	db *sqlx.DB
}

// VehicleIssuedRepoDI ...
type VehicleIssuedRepoDI struct {
	dig.In
	DB *sqlx.DB
}

// NewVehicleIssuedRepo ...
func NewVehicleIssuedRepo(param VehicleIssuedRepoDI) interfaces.VehicleIssuedRepo {
	return &VehicleIssuedRepo{
		db: param.DB,
	}
}

// FindByRegistryNumber ...
func (r *VehicleIssuedRepo) FindByRegistryNumber(registryNumber string) (*entity.Issued, error) {
	var record entity.Issued

	err := r.db.Get(&record, "SELECT "+
		"id, created_at, created_by, company_inn, company_ogrn, company_name, "+
		"company_fio, company_car, legal_basement, pass_number, district, pass_type, issued_at, "+
		"registry_number, file_id, shipping "+
		"FROM issued where registry_number = ?", registryNumber)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "unable get issued record by registry_number %s", registryNumber)
	}

	return &record, nil
}

// FindByPass ...
func (r *VehicleIssuedRepo) FindByPass(pass string) (*entity.Issued, error) {
	var record entity.Issued

	err := r.db.Get(&record, "SELECT "+
		"id, created_at, created_by, company_inn, company_ogrn, company_name, "+
		"company_fio, company_car, legal_basement, pass_number, district, pass_type, issued_at, "+
		"registry_number, file_id, shipping "+
		"FROM issued where pass_number = ?", pass)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "unable get issued record by pass_number %s", pass)
	}

	return &record, nil
}

// FindByCar ...
func (r *VehicleIssuedRepo) FindByCar(car string) (*entity.Issued, error) {
	var record entity.Issued

	query := "SELECT " +
		"id, created_at, created_by, company_inn, company_ogrn, company_name, " +
		"company_fio, company_car, legal_basement, pass_number, district, pass_type, issued_at, " +
		"registry_number, file_id, shipping " +
		" FROM issued where company_car = ?"

	err := r.db.Get(&record, query, car)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "unable get issued record by company_car %s", car)
	}

	return &record, nil
}

// Create ...
func (r *VehicleIssuedRepo) Create(data *entity.Issued) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		query := "INSERT INTO issued (" +
			"created_at, created_by, " +
			"company_inn, company_ogrn, company_name, company_fio, company_car, " +
			"legal_basement, pass_number, pass_type, " +
			"district, issued_at, registry_number, shipping, file_id" +
			") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

		res, err := t.Exec(query,
			data.CreatedAt,
			data.CreatedBy,
			data.CompanyInn,
			data.CompanyOgrn,
			data.CompanyName,
			data.CompanyFio,
			data.CompanyCar,
			data.LegalBasement,
			data.PassNumber,
			data.PassType,
			data.District,
			data.IssuedAt,
			data.RegistryNumber,
			data.Shipping,
			data.FileID)

		if err != nil {
			return err
		}

		data.ID, err = res.LastInsertId()
		if err != nil {
			return errors.Wrap(err, "unable get bid bids lastInsertID")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}

// Update ...
func (r *VehicleIssuedRepo) Update(data *entity.Issued) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "UPDATE issued SET " +
			"created_at = ?, created_by = ?, " +
			"company_inn = ?, company_ogrn = ?, company_name = ?, company_fio = ?, company_car = ?, " +
			"legal_basement = ?, pass_number = ?, district = ?, " +
			"pass_type = ?, issued_at = ?, registry_number = ?, shipping = ?, file_id = ? " +
			"WHERE id = ?"
		_, err := t.Exec(sql,
			data.CreatedAt,
			data.CreatedBy,
			data.CompanyInn,
			data.CompanyOgrn,
			data.CompanyName,
			data.CompanyFio,
			data.CompanyCar,
			data.LegalBasement,
			data.PassNumber,
			data.District,
			data.PassType,
			data.IssuedAt,
			data.RegistryNumber,
			data.Shipping,
			data.FileID,
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
func (r *VehicleIssuedRepo) Read(id int64) (*entity.Issued, error) {
	var data entity.Issued

	err := r.db.Get(&data, "SELECT "+
		"created_at, created_by, company_inn, company_ogrn, company_name, "+
		"company_fio, company_car, legal_basement, pass_number, district, pass_type, issued_at, "+
		"registry_number, file_id, shipping "+
		"FROM issued WHERE id = ?", id)
	if err != nil {
		return nil, errors.Wrapf(err, "unable get issued record by id %s", id)
	}

	return &data, nil
}

// Delete ...
func (r *VehicleIssuedRepo) Delete(id int64) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "DELETE FROM issued WHERE id = ?"
		_, err := t.Exec(sql, id)
		if err != nil {
			return errors.Wrapf(err, "unable delete from issued by id %d", id)
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}
