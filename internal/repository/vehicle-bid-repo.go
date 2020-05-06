package repository

import (
	"github.com/alexey-zayats/claim-parser/internal/database"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// VehicleBidRepo ...
type VehicleBidRepo struct {
	db *sqlx.DB
}

// VehicleBidRepoDI ...
type VehicleBidRepoDI struct {
	dig.In
	DB *sqlx.DB
}

// NewVehicleBidRepo ...
func NewVehicleBidRepo(di VehicleBidRepoDI) interfaces.VehicleBidRepo {
	return &VehicleBidRepo{
		db: di.DB,
	}
}

// Create ...
func (r *VehicleBidRepo) Create(data *entity.Bid) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		query := "INSERT INTO bids (" +
			"file_id, company_id, branch_id, company_branch, " +
			"company_name, company_address, company_ceo_phone, " +
			"company_ceo_email, company_ceo_name," +
			"agree, confirm, " +
			"workflow_status, code, district_id, pass_type, " +
			"created_at, created_by, source, " + // user_id,
			"city, who_address, phone_where, fio_where, city_where, " +
			"address_where, `from`, `to`, reason, link" +
			") VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)" // , ?

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
			data.WorkflowStatus,
			data.Agree,
			data.Confirm,
			data.Code,
			data.DistrictID,
			data.PassType,
			data.CreatedAt,
			data.CreatedBy,
			//data.UserID,
			data.Source,

			data.CityFrom,
			data.AddressDest,
			data.WhoNeedsHelpPhone,
			data.WhoNeedsHelp,
			data.CityTo,
			data.AddressWhere,
			data.DateFrom,
			data.DateTo,
			data.OtherReason,
			data.DocLinks,
		)

		if err != nil {
			return errors.Wrap(err, "unable create bid")
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
func (r *VehicleBidRepo) Update(data *entity.Bid) error {

	err := database.WithTransaction(r.db, func(t database.Transaction) error {

		sql := "UPDATE bids SET " +
			"file_id = ?, company_id = ?, branch_id = ?, company_branch = ?, " +
			"company_name = ?, company_address = ?, company_ceo_phone = ?, " +
			"company_ceo_email = ?, company_ceo_name = ?," +
			"agree = ?, confirm = ?, " +
			"workflow_status = ?, code = ?, district_id = ?, pass_type = ?, " +
			"created_at = ?, created_by = ?, source = ?, " + // user_id = ?,
			"city = ?, who_address = ?, phone_where = ?, fio_where = ?, city_where = ?, " +
			"address_where = ?, `from` = ?, `to` = ?, reason = ?, link = ?" +
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
			data.WorkflowStatus,
			data.Agree,
			data.Confirm,
			data.Code,
			data.DistrictID,
			data.PassType,
			data.CreatedAt,
			data.CreatedBy,
			//data.UserID,
			data.Source,

			data.CityFrom,
			data.AddressDest,
			data.WhoNeedsHelpPhone,
			data.WhoNeedsHelp,
			data.CityTo,
			data.AddressWhere,
			data.DateFrom,
			data.DateTo,
			data.OtherReason,
			data.DocLinks,

			data.ID)

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
func (r *VehicleBidRepo) Read(id int64) (*entity.Bid, error) {
	var request entity.Bid

	err := r.db.Get(&request, "SELECT "+
		"id, "+
		"file_id, company_id, branch_id, company_branch, "+
		"company_name, company_address, company_ceo_phone, "+
		"company_ceo_email, company_ceo_name,"+
		"agree, confirm, "+
		"workflow_status, code, district_id, pass_type, "+
		"created_at, created_by, user_id, source, "+
		"city, who_address, phone_where, fio_where, city_where, "+
		"address_where, `from`, `to`, reason, link "+
		"FROM bids WHERE id = ?", id)
	if err != nil {
		return nil, errors.Wrapf(err, "unable get bids record by id %s", id)
	}

	return &request, nil
}

// Delete ...
func (r *VehicleBidRepo) Delete(id int64) error {

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
