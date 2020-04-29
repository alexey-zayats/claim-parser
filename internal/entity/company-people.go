package entity

// CompanyPeople - компании по пропускам для сотрудников
type CompanyPeople struct {
	ID       int64  `db:"id"`
	OGRN     string `db:"ogrn"`
	INN      string `db:"inn"`
	Name     string `db:"name"`
	BranchID int64  `db:"branch_id"`
	Status   int    `db:"status"`
}
