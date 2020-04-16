package model

// Event ...
type Event struct {
	CreatedBy  int64  `json:"created_by"`
	Filepath   string `json:"filepath"`
	Source     int    `json:"source"`
	PassType   int    `json:"type"`
	DistrictID int64  `json:"district"`
	FileID     int64  `json:"file_id"`
	Check      int    `json:"check"`
	BranchID   int64  `json:"branch_id"`
}
