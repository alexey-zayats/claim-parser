package model

// Event ...
type Event struct {
	CreatedBy int    `json:"created_by"`
	Filepath  string `json:"filepath"`
	Source    int    `json:"source"`
	PassType  int    `json:"type"`
	District  int    `json:"district"`
	FileID    int    `json:"file_id"`
}

// State ...
type State struct {
	ID     int
	Status int
	Error  error `json:"omitempty"`
}
