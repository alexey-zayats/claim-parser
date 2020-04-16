package interfaces

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
)

// FileRepository ...
type FileRepository interface {
	Create(data *entity.File) error
	Read(id int64) (*entity.File, error)
	UpdateState(data *entity.File) error
	Delete(id int64) error
}
