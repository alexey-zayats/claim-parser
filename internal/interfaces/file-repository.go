package interfaces

import "github.com/alexey-zayats/claim-parser/internal/model"

// FileRepository ...
type FileRepository interface {
	Create(data *model.File) error
	Read(id int64) (*model.File, error)
	UpdateState(data *model.File) error
	Delete(id int64) error
}
