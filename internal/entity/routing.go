package entity

// Routing ...
type Routing struct {
	ID         int64 `db:"id"`
	SourceID   int64 `db:"source_id"`
	DistrictID int64 `db:"district_id"`
	CleanID    int64 `db:"clean_id"`
	DirtyID    int64 `db:"dirty_id"`
}
