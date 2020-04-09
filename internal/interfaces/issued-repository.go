package interfaces

import "github.com/alexey-zayats/claim-parser/internal/model"

// IssuedRepository ...
type IssuedRepository interface {
	Create(data *model.Issued) error
	Read(id int64) (*model.Issued, error)
	Update(data *model.Issued) error
	Delete(id int64) error

	FindByCar(car string) (*model.Issued, error)
	FindByPass(pass string) (*model.Issued, error)
}
