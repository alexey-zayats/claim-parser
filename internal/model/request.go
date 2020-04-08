package model

import (
	"time"
)

// Request ...
type Request struct {
	ID             int       `db:"id"`
	FileID         int       `db:"file_id"`
	Status         int       `db:"status"`
	WorkflowStatus int       `db:"workflow_status"`
	Code           string    `db:"code"`
	District       int       `db:"district"`
	PassType       int       `db:"type"`
	CreatedAt      time.Time `db:"created_at"`
	CreatedBy      int       `db:"created_by"`
	UserID         int       `db:"user_id"`
	Source         string    `db:"source"`
}
