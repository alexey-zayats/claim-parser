package entity

import (
	"time"
)

// Pass ...
type Pass struct {
	ID         int64     `db:"id"`
	BidID      int64     `db:"bid_id"`
	IssuedID   int64     `db:"issued_id"`
	Lastname   string    `db:"lastname"`
	Firstname  string    `db:"firstname"`
	Patrname   string    `db:"patrname"`
	Car        string    `db:"car"`
	Source     int64     `db:"source"`
	DistrictID int64     `db:"district_id"`
	PassType   int       `db:"pass_type"`
	PassNumber string    `db:"pass_number"`
	Shipping   int       `db:"shipping"`
	Status     int       `db:"status"`
	FileID     int64     `db:"file_id"`
	CreatedAt  time.Time `db:"created_at"`
	CreatedBy  int64     `db:"created_by"`
}
