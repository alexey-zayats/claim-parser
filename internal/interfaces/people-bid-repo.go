package interfaces

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
)

// PeopleBidRepo ...
type PeopleBidRepo interface {
	Create(*entity.BidPeople) error
	Read(id int64) (*entity.BidPeople, error)
	Update(*entity.BidPeople) error
	Delete(id int64) error
}
