package interfaces

import "github.com/alexey-zayats/claim-parser/internal/model"

// BidRepository ...
type BidRepository interface {
	Create(*model.Bid) error
	Read(id int64) (*model.Bid, error)
	Update(*model.Bid) error
	Delete(id int64) error
}
