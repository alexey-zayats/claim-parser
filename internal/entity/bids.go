package entity

import (
	"time"
)

// Bid ...
type Bid struct {
	ID                int64     `db:"id"`
	FileID            int64     `db:"file_id"`
	CompanyID         int64     `db:"company_id"`
	BranchID          int64     `db:"branch_id"`
	CompanyBranch     string    `db:"company_branch"`
	CompanyName       string    `db:"company_name"`
	CompanyAddress    string    `db:"company_address"` // address_where
	CompanyCeoPhone   string    `db:"company_ceo_phone"`
	CompanyCeoEmail   string    `db:"company_ceo_email"`
	CompanyCeoName    string    `db:"company_ceo_name"`
	Agree             int       `db:"agree"`
	Confirm           int       `db:"confirm"`
	WorkflowStatus    int       `db:"workflow_status"`
	Code              string    `db:"code"`
	DistrictID        int64     `db:"district_id"`
	PassType          int       `db:"pass_type"`
	CreatedAt         time.Time `db:"created_at"`
	CreatedBy         int64     `db:"created_by"`
	UserID            int64     `db:"user_id"`
	Source            string    `db:"source"`
	MovedTo           int64     `db:"moved_to"`
	AlighnedID        int64     `db:"alighned_id"`
	PrintID           int64     `db:"print_id"`
	CityFrom          string    `db:"city"`
	CityTo            string    `db:"city_where"`
	AddressDest       string    `db:"who_address"`
	AddressWhere      string    `db:"address_where"`
	WhoNeedsHelpPhone string    `db:"phone_where"`
	WhoNeedsHelp      string    `db:"fio_where"`
	DateFrom          time.Time `db:"from"`
	DateTo            time.Time `db:"to"`
	OtherReason       string    `db:"reason"`
	DocLinks          string    `db:"link"`
}
