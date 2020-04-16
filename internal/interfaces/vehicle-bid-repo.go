package interfaces

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
)

// VehicleBidRepo ...
type VehicleBidRepo interface {
	Create(*entity.Bid) error
	Read(id int64) (*entity.Bid, error)
	Update(*entity.Bid) error
	Delete(id int64) error
}
