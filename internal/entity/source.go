package entity

// Source ...
type Source struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Title string `db:"title"`
}
