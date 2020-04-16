package entity

import "time"

// File ...
type File struct {
	ID        int64     `db:"id"`
	Filepath  string    `db:"filepath"`
	Status    int       `db:"status"`
	Log       string    `db:"log"`
	Source    string    `db:"source"`
	CreatedAt time.Time `db:"created_at"`
}
