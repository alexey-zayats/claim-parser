package repository

import (
	"github.com/alexey-zayats/claim-parser/internal/database"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// PeoplePassRepo ...
type PeoplePassRepo struct {
	db *sqlx.DB
}

// PeoplePassRepoDI ...
type PeoplePassRepoDI struct {
	dig.In
	DB *sqlx.DB
}

// NewPeoplePassRepo ...
func NewPeoplePassRepo(di PeoplePassRepoDI) interfaces.PeoplePassRepo {
	return &PeoplePassRepo{
		db: di.DB,
	}
}

// Create ...
func (r *PeoplePassRepo) Create(data *entity.PassPeople) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		query := "INSERT INTO passes_people (" +
			"bid_id, source, " +
			"district_id, pass_type, pass_number, shipping, " +
			"status, lastname, firstname, patrname " +
			") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

		res, err := t.Exec(query,
			data.BidID,
			data.Source,
			data.DistrictID,
			data.PassType,
			data.PassNumber,
			data.Shipping,
			data.Status,
			data.Lastname,
			data.Firstname,
			data.Patrname,
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
func (r *PeoplePassRepo) Update(data *entity.PassPeople) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "UPDATE passes_people SET " +
			"bid_id = ?, issued_id = ?, source = ?, " +
			"district_id = ?, pass_type = ?, pass_number = ?, " +
			"shipping = ?, status = ?, " +
			"lastname = ?, firstname = ?, patrname = ? " +
			"WHERE id = ?"

		_, err := t.Exec(sql,
			data.BidID,
			data.IssuedID,
			data.Source,
			data.DistrictID,
			data.PassType,
			data.PassNumber,
			data.Shipping,
			data.Status,
			data.Lastname,
			data.Firstname,
			data.Patrname,
			// --
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
func (r *PeoplePassRepo) Read(id int64) (*entity.PassPeople, error) {
	var pass entity.PassPeople

	err := r.db.Get(&pass, "SELECT "+
		"id, "+
		"bid_id, issued_id, source, "+
		"district_id, pass_type, pass_number, shipping, "+
		"status, lastname, firstname, patrname "+
		"FROM passes_people WHERE id = ?", id)
	if err != nil {
		return nil, errors.Wrapf(err, "unable get passes_people record by id %s", id)
	}

	return &pass, nil
}

// Delete ...
func (r *PeoplePassRepo) Delete(id int64) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "DELETE FROM passes_people WHERE id = ?"
		_, err := t.Exec(sql, id)
		if err != nil {
			return errors.Wrapf(err, "unable delete from passes_people by id %d", id)
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}
