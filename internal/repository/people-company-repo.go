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

// PeopleCompanyRepo ...
type PeopleCompanyRepo struct {
	db *sqlx.DB
}

// PeopleCompanyRepoInput ...
type PeopleCompanyRepoInput struct {
	dig.In
	DB *sqlx.DB
}

// NewPeopleCompanyRepo ...
func NewPeopleCompanyRepo(param PeopleCompanyRepoInput) interfaces.PeopleCompanyRepo {
	return &PeopleCompanyRepo{
		db: param.DB,
	}
}

// FindByOgrnInn ...
func (r *PeopleCompanyRepo) FindByOgrnInn(ogrn, inn int64) (*entity.CompanyPeople, error) {
	var record entity.CompanyPeople

	query :=
		"SELECT id, ogrn, inn, name, branch_id, status " +
			"FROM companies_people " +
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
func (r *PeopleCompanyRepo) FindByINN(inn int64) (*entity.CompanyPeople, error) {
	var record entity.CompanyPeople

	query := "SELECT id, ogrn, inn, name, branch_id, status " +
		"FROM companies_people WHERE inn = ?"

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
func (r *PeopleCompanyRepo) FindByOGRN(ogrn int64) (*entity.CompanyPeople, error) {
	var record entity.CompanyPeople

	query := "SELECT " +
		"id, ogrn, inn, name, branch_id, status " +
		"FROM " +
		"companies_people " +
		"WHERE " +
		"ogrn = ?"

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
func (r *PeopleCompanyRepo) Create(data *entity.CompanyPeople) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		query := "INSERT INTO companies_people (" +
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
func (r *PeopleCompanyRepo) Update(data *entity.CompanyPeople) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "UPDATE companies_people " +
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
func (r *PeopleCompanyRepo) Read(id int64) (*entity.CompanyPeople, error) {
	var pass entity.CompanyPeople

	query :=
		"SELECT " +
			"id, " +
			"ogrn, inn, name, branch_id, status " +
			"FROM " +
			"companies_people " +
			"WHERE " +
			"id = ?"

	err := r.db.Get(&pass, query, id)
	if err != nil {
		return nil, errors.Wrapf(err, "unable get companies_people record by id %s", id)
	}

	return &pass, nil
}

// Delete ...
func (r *PeopleCompanyRepo) Delete(id int64) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "DELETE FROM passes WHERE id = ?"
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
