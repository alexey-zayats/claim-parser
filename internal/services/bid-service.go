package services

import (
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"go.uber.org/dig"
)

// BidService ...
type BidService struct {
	repo interfaces.BidRepository
}

// BidServiceDI ...
type BidServiceDI struct {
	dig.In
	Repo interfaces.BidRepository
}

// NewBidService ...
func NewBidService(di BidServiceDI) *BidService {
	return &BidService{
		repo: di.Repo,
	}
}

// Create ...
func (s *BidService) Create(data *model.Bid) error {
	return s.repo.Create(data)
}

// Update ...
func (s *BidService) Update(data *model.Bid) error {
	return s.repo.Update(data)
}

// Delete ...
func (s *BidService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Read ...
func (s *BidService) Read(id int64) (*model.Bid, error) {
	return s.repo.Read(id)
}
