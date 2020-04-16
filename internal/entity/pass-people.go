package entity

// PassPeople ...
type PassPeople struct {
	ID         int64  `db:"id"`
	BidID      int64  `db:"bid_id"`
	IssuedID   int64  `db:"issued_id"`
	Source     int    `db:"source"`
	DistrictID int64  `db:"district_id"`
	PassType   int    `db:"pass_type"`
	PassNumber string `db:"pass_number"`
	Shipping   int    `db:"shipping"`
	Status     int    `db:"status"`
	Lastname   string `db:"lastname"`
	Firstname  string `db:"firstname"`
	Patrname   string `db:"patrname"`
}
