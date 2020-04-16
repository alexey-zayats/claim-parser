package services

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"go.uber.org/dig"
)

// PeopleBidService ...
type PeopleBidService struct {
	repo interfaces.PeopleBidRepo
}

// PeopleBidServiceDI ...
type PeopleBidServiceDI struct {
	dig.In
	Repo interfaces.PeopleBidRepo
}

// NewPeopleBidService ...
func NewPeopleBidService(di PeopleBidServiceDI) *PeopleBidService {
	return &PeopleBidService{
		repo: di.Repo,
	}
}

// Create ...
func (s *PeopleBidService) Create(data *entity.BidPeople) error {
	return s.repo.Create(data)
}

// Update ...
func (s *PeopleBidService) Update(data *entity.BidPeople) error {
	return s.repo.Update(data)
}

// Delete ...
func (s *PeopleBidService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Read ...
func (s *PeopleBidService) Read(id int64) (*entity.BidPeople, error) {
	return s.repo.Read(id)
}
