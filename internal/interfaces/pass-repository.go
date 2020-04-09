package interfaces

import "github.com/alexey-zayats/claim-parser/internal/model"

// PassRepository ...
type PassRepository interface {
	Create(*model.Pass) error
	Read(id int64) (*model.Pass, error)
	Update(*model.Pass) error
	Delete(id int64) error
}
