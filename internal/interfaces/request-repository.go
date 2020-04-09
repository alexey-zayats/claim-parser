package interfaces

import "github.com/alexey-zayats/claim-parser/internal/model"

// BidRepository ...
type BidRepository interface {
	Create(*model.Bid) (int64, error)
	Read(id int) (*model.Bid, error)
	Update(*model.Bid) error
	Delete(id int) error
}
