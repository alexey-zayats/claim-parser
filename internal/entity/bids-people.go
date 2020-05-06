package entity

import "time"

// BidPeople ...
type BidPeople struct {
	ID              int64     `db:"ID"`
	FileID          int64     `db:"file_id"`
	CompanyID       int64     `db:"company_id"`
	BranchID        int64     `db:"branch_id"`
	CompanyBranch   string    `db:"company_branch"`
	CompanyName     string    `db:"company_name"`
	CompanyAddress  string    `db:"company_address"`
	CompanyCeoPhone string    `db:"company_ceo_phone"`
	CompanyCeoEmail string    `db:"company_ceo_email"`
	CompanyCeoName  string    `db:"company_ceo_name"`
	Agree           int       `db:"agree"`
	Confirm         int       `db:"confirm"`
	WorkflowStatus  int       `db:"workflow_status"`
	DistrictID      int64     `db:"district_id"`
	PassType        int       `db:"pass_type"`
	Source          string    `db:"source"`
	UserID          int64     `db:"user_id"`
	MovedTo         int64     `db:"moved_to"`
	AlighnedID      int       `db:"alighned_id"`
	PrintID         int       `db:"print_id"`
	CreatedAt       time.Time `db:"created_at"`
	CreatedBy       int64     `db:"created_by"`
}
