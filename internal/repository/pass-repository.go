package repository

import (
	"github.com/alexey-zayats/claim-parser/internal/database"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// PassRepository ...
type PassRepository struct {
	db *sqlx.DB
}

// PassRepositoryInput ...
type PassRepositoryInput struct {
	dig.In
	DB *sqlx.DB
}

// NewPassRepository ...
func NewPassRepository(param PassRepositoryInput) interfaces.PassRepository {
	return &PassRepository{
		db: param.DB,
	}
}

// Create ...
func (r *PassRepository) Create(data *model.Pass) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		query := "INSERT INTO passes (" +
			"company_branch, company_okved, company_inn, company_name, company_address, company_ceo_phone," +
			"company_ceo_email, company_lastname, company_firstname, company_patrname, " +
			"employee_lastname, employee_firstname, employee_patrname, employee_car, employee_agree, employee_confirm, " +
			"source, district, type, number, status, file_id, created_at, created_by, bid_id" +
			") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

		res, err := t.Exec(query,
			data.CompanyBranch,
			data.CompanyOkved,
			data.CompanyInn,
			data.CompanyName,
			data.CompanyAddress,
			data.CompanyCeoPhone,
			data.CompanyCeoEmail,
			data.CompanyLastname,
			data.CompanyFirstname,
			data.CompanyPatrname,
			data.EmployeeLastname,
			data.EmployeeFirstname,
			data.EmployeePatrname,
			data.EmployeeCar,
			data.EmployeeAgree,
			data.EmployeeConfirm,
			data.Source,
			data.District,
			data.PassType,
			data.PassNumber,
			data.Status,
			data.FileID,
			data.CreatedAt,
			data.CreatedBy,
			data.BidID)

		if err != nil {
			return errors.Wrap(err, "unable update passes")
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
func (r *PassRepository) Update(data *model.Pass) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "UPDATE passes SET " +
			"company_branch = ?, company_okved = ?, company_inn = ?, company_name = ?, company_address = ?, " +
			"company_ceo_phone = ?, company_ceo_email = ?, company_lastname = ?, company_firstname = ?, " +
			"company_patrname = ?, employee_lastname = ?, employee_firstname = ?, employee_patrname = ?, " +
			"employee_car = ?, employee_agree = ?, employee_confirm = ?, source = ?, district = ?, " +
			"type = ?, number = ?, status = ?, file_id = ?, created_at = ?, created_by = ?, bid_id = ?" +
			"WHERE id = ?"

		_, err := t.Exec(sql,
			data.CompanyBranch,
			data.CompanyOkved,
			data.CompanyInn,
			data.CompanyName,
			data.CompanyAddress,
			data.CompanyCeoPhone,
			data.CompanyCeoEmail,
			data.CompanyLastname,
			data.CompanyFirstname,
			data.CompanyPatrname,
			data.EmployeeLastname,
			data.EmployeeFirstname,
			data.EmployeePatrname,
			data.EmployeeCar,
			data.EmployeeAgree,
			data.EmployeeConfirm,
			data.Source,
			data.District,
			data.PassType,
			data.PassNumber,
			data.Status,
			data.FileID,
			data.CreatedAt,
			data.CreatedBy,
			data.ID,
			data.BidID)

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
func (r *PassRepository) Read(id int64) (*model.Pass, error) {
	var pass *model.Pass

	err := r.db.Get(pass, "select * from passes where id=?", id)
	if err != nil {
		return nil, errors.Wrapf(err, "unable get passes record by id %s", id)
	}

	return pass, nil
}

// Delete ...
func (r *PassRepository) Delete(id int64) error {

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
