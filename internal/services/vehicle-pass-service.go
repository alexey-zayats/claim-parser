package services

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"go.uber.org/dig"
)

// VehiclePassService ...
type VehiclePassService struct {
	repo interfaces.VehiclePassRepo
}

// VehiclePassServiceDI ...
type VehiclePassServiceDI struct {
	dig.In
	Repo interfaces.VehiclePassRepo
}

// NewVehiclePassService ...
func NewVehiclePassService(di VehiclePassServiceDI) *VehiclePassService {
	return &VehiclePassService{
		repo: di.Repo,
	}
}

// Create ...
func (s *VehiclePassService) Create(data *entity.Pass) error {
	return s.repo.Create(data)
}

// Update ...
func (s *VehiclePassService) Update(data *entity.Pass) error {
	return s.repo.Update(data)
}

// Delete ...
func (s *VehiclePassService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Read ...
func (s *VehiclePassService) Read(id int64) (*entity.Pass, error) {
	return s.repo.Read(id)
}

// FindByCar ...
func (s *VehiclePassService) FindByCar(car string) (*entity.Pass, error) {
	return s.repo.FindByCar(car)
}
