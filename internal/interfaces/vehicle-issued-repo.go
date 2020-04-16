package interfaces

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
)

// VehicleIssuedRepo ...
type VehicleIssuedRepo interface {
	Create(data *entity.Issued) error
	Read(id int64) (*entity.Issued, error)
	Update(data *entity.Issued) error
	Delete(id int64) error

	FindByCar(car string) (*entity.Issued, error)
	FindByPass(pass string) (*entity.Issued, error)
}
