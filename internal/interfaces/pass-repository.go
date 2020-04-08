package interfaces

import "github.com/alexey-zayats/claim-parser/internal/model"

// PassRepository ...
type PassRepository interface {
	Create(*model.Pass) (int64, error)
	Read(id int) (*model.Pass, error)
	Update(*model.Pass) error
	Delete(id int) error
}
