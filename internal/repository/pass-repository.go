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

	logrus.WithFields(logrus.Fields{"data": data}).Debug("PassRepository.Create")

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		query := "INSERT INTO passes (" +
			"company_branch, company_okved, company_inn, company_name, company_address, company_ceo_phone," +
			"company_ceo_email, company_lastname, company_firstname, company_patrname, " +
			"employee_lastname, employee_firstname, employee_patrname, employee_car, employee_agree, employee_confirm, " +
			"source, district, type, number, status, file_id, created_at, created_by" +
			") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

		_, err := t.Exec(query,
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
			data.CreatedBy)

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
