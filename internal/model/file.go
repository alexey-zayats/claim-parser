package model

import "time"

// File ...
type File struct {
	ID        int64     `db:"id"`
	Filepath  string    `db:"filepath"`
	Status    int       `db:"status"`
	Log       string    `db:"log"`
	CreatedAt time.Time `db:"created_at"`
	Source    string    `db:"source"`
}
