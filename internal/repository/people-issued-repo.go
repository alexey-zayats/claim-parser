package repository

import (
	"github.com/alexey-zayats/claim-parser/internal/database"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// PeopleIssuedRepo ...
type PeopleIssuedRepo struct {
	db *sqlx.DB
}

// PeopleIssuedRepoDI ...
type PeopleIssuedRepoDI struct {
	dig.In
	DB *sqlx.DB
}

// NewPeopleIssuedRepo ...
func NewPeopleIssuedRepo(param PeopleIssuedRepoDI) interfaces.PeopleIssuedRepo {
	return &PeopleIssuedRepo{
		db: param.DB,
	}
}

// Create ...
func (r *PeopleIssuedRepo) Create(data *entity.IssuedPeople) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		query := "INSERT INTO issued_people " +
			"(" +
			"district_id, company_id, lastname, firstname, patrname, " +
			"legal_basement, pass_number, created_at, created_by, issued_at, " +
			"registry_number, shiping, arm_number, arm_number_by, arm_number_at" +
			") " +
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

		res, err := t.Exec(query,
			data.DistrictID,
			data.CompanyID,
			data.Lastname,
			data.Firstname,
			data.Patrname,
			data.LegalBasement,
			data.PassNumber,
			data.CreatedAt,
			data.CreatedBy,
			data.IssuedAt,
			data.RegistryNumber,
			data.Shiping,
			data.ArmNumber,
			data.ArmNumberBy,
			data.ArmNumberAt,
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
func (r *PeopleIssuedRepo) Update(data *entity.IssuedPeople) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "UPDATE issued_people SET " +
			"district_id = ?, company_id = ?, lastname = ?, firstname = ?, patrname = ?, " +
			"legal_basement = ?, pass_number = ?, created_at = ?, created_by = ?, issued_at = ?, " +
			"registry_number = ?, shiping = ?, arm_number = ?, arm_number_by = ?, arm_number_at = ?" +
			"WHERE id = ?"
		_, err := t.Exec(sql,
			data.DistrictID,
			data.CompanyID,
			data.Lastname,
			data.Firstname,
			data.Patrname,
			data.LegalBasement,
			data.PassNumber,
			data.CreatedAt,
			data.CreatedBy,
			data.IssuedAt,
			data.RegistryNumber,
			data.Shiping,
			data.ArmNumber,
			data.ArmNumberBy,
			data.ArmNumberAt,
			// ---
			data.ID,
		)

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
func (r *PeopleIssuedRepo) Read(id int64) (*entity.IssuedPeople, error) {
	var data entity.IssuedPeople

	query := "SELECT " +
		"district_id, company_id, lastname, firstname, patrname, " +
		"legal_basement, pass_number, created_at, created_by, issued_at, " +
		"registry_number, shiping, arm_number, arm_number_by, arm_number_at " +
		"FROM issued_people WHERE id = ?"

	err := r.db.Get(&data, query, id)
	if err != nil {
		return nil, errors.Wrapf(err, "unable get issued_people record by id %s", id)
	}

	return &data, nil
}

// Delete ...
func (r *PeopleIssuedRepo) Delete(id int64) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "DELETE FROM issued WHERE id = ?"
		_, err := t.Exec(sql, id)
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
