package repository

import (
	"github.com/alexey-zayats/claim-parser/internal/database"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// PeopleBidRepo ...
type PeopleBidRepo struct {
	db *sqlx.DB
}

// PeopleBidRepoDI ...
type PeopleBidRepoDI struct {
	dig.In
	DB *sqlx.DB
}

// NewPeopleBidRepo ...
func NewPeopleBidRepo(di PeopleBidRepoDI) interfaces.PeopleBidRepo {
	return &PeopleBidRepo{
		db: di.DB,
	}
}

// Create ...
func (r *PeopleBidRepo) Create(data *entity.BidPeople) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		query := "INSERT INTO bids_people (" +
			"file_id, company_id, branch_id, " +
			"company_branch, company_name, " +
			"company_address, company_ceo_phone, company_ceo_email, " +
			"company_ceo_name, agree, confirm, workflow_status, district_id, pass_type, source, " +
			"user_id, moved_to, alighned_id, print_id, created_at, created_by" +
			") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"

		res, err := t.Exec(query,
			data.FileID,
			data.CompanyID,
			data.BranchID,
			data.CompanyBranch,
			data.CompanyName,
			data.CompanyAddress,
			data.CompanyCeoPhone,
			data.CompanyCeoEmail,
			data.CompanyCeoName,
			data.Agree,
			data.Confirm,
			data.WorkflowStatus,
			data.DistrictID,
			data.PassType,
			data.Source,
			data.UserID,
			data.MovedTo,
			data.AlighnedID,
			data.PrintID,
			data.CreatedAt,
			data.CreatedBy,
		)

		if err != nil {
			return errors.Wrap(err, "unable create bids_people")
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
func (r *PeopleBidRepo) Update(data *entity.BidPeople) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "UPDATE bids_people SET " +
			"file_id = ?, company_id = ?, branch_id = ?, " +
			"company_branch = ?, company_name = ?, " +
			"company_address = ?, company_ceo_phone = ?, company_ceo_email = ?, " +
			"company_ceo_name = ?, agree = ?, confirm = ?, workflow_status = ?, district_id = ?, pass_type = ?, source = ?, " +
			"user_id = ?, moved_to = ?, alighned_id = ?, print_id = ?, created_at = ?, created_by = ? " +
			"WHERE id = ?"
		_, err := t.Exec(sql,
			data.FileID,
			data.CompanyID,
			data.BranchID,
			data.CompanyBranch,
			data.CompanyName,
			data.CompanyAddress,
			data.CompanyCeoPhone,
			data.CompanyCeoEmail,
			data.CompanyCeoName,
			data.Agree,
			data.Confirm,
			data.WorkflowStatus,
			data.DistrictID,
			data.PassType,
			data.Source,
			data.UserID,
			data.MovedTo,
			data.AlighnedID,
			data.PrintID,
			data.CreatedAt,
			data.CreatedBy,
			// ----
			data.ID,
		)

		if err != nil {
			return errors.Wrap(err, "unable update bids_people")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}

// Read ...
func (r *PeopleBidRepo) Read(id int64) (*entity.BidPeople, error) {
	var request entity.BidPeople

	query := "SELECT " +
		"file_id, company_id, branch_id, " +
		"company_branch, company_name, " +
		"company_address, company_ceo_phone, company_ceo_email, " +
		"company_ceo_name, agree, confirm, workflow_status, district_id, pass_type, source, " +
		"user_id, moved_to, alighned_id, print_id, created_at, created_by " +
		"FROM bids_people " +
		"WHERE id = ?"

	err := r.db.Get(&request, query, id)
	if err != nil {
		return nil, errors.Wrapf(err, "unable get bids_people record by id %s", id)
	}

	return &request, nil
}

// Delete ...
func (r *PeopleBidRepo) Delete(id int64) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "DELETE FROM bids_people WHERE id = ?"
		_, err := t.Exec(sql, id)
		if err != nil {
			return errors.Wrapf(err, "unable delete from bids_people by id %d", id)
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "transaction error")
	}

	return nil
}
