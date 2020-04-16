package entity

// Branch ...
type Branch struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
	Type string `db:"type"`
}
