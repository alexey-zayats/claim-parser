package interfaces

import "github.com/alexey-zayats/claim-parser/internal/model"

// FileRepository ...
type FileRepository interface {
	Create(data *model.File) (int64, error)
	Read(id int) (*model.File, error)
	Update(data *model.File) error
	Delete(id int) error
}
