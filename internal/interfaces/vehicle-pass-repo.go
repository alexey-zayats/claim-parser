package interfaces

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
)

// VehiclePassRepo ...
type VehiclePassRepo interface {
	Create(*entity.Pass) error
	Read(id int64) (*entity.Pass, error)
	Update(*entity.Pass) error
	Delete(id int64) error

	FindByCar(car string) (*entity.Pass, error)
}
