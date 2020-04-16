package entity

// Company - компании по пропускам для транспорта
type Company struct {
	ID       int64  `db:"id"`
	OGRN     int64  `db:"ogrn"`
	INN      int64  `db:"inn"`
	Name     string `db:"name"`
	BranchID int64  `db:"branch_id"`
	Status   int    `db:"status"`
}
