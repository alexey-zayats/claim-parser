package services

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"go.uber.org/dig"
)

// VehicleBidService ...
type VehicleBidService struct {
	repo interfaces.VehicleBidRepo
}

// VehicleBidServiceDI ...
type VehicleBidServiceDI struct {
	dig.In
	Repo interfaces.VehicleBidRepo
}

// NewVehicleBidService ...
func NewVehicleBidService(di VehicleBidServiceDI) *VehicleBidService {
	return &VehicleBidService{
		repo: di.Repo,
	}
}

// Create ...
func (s *VehicleBidService) Create(data *entity.Bid) error {
	return s.repo.Create(data)
}

// Update ...
func (s *VehicleBidService) Update(data *entity.Bid) error {
	return s.repo.Update(data)
}

// Delete ...
func (s *VehicleBidService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Read ...
func (s *VehicleBidService) Read(id int64) (*entity.Bid, error) {
	return s.repo.Read(id)
}
