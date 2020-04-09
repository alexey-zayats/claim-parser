package model

import (
	"time"
)

// Bid ...
type Bid struct {
	ID             int64     `db:"id"`
	FileID         int64     `db:"file_id"`
	WorkflowStatus int       `db:"workflow_status"`
	Code           string    `db:"code"`
	District       int64     `db:"district"`
	PassType       int       `db:"type"`
	CreatedAt      time.Time `db:"created_at"`
	CreatedBy      int64     `db:"created_by"`
	UserID         int64     `db:"user_id"`
	Source         string    `db:"source"`
}
