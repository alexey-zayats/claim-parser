package services

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"go.uber.org/dig"
)

// VehicleCompanyService ...
type VehicleCompanyService struct {
	repo interfaces.VehicleCompanyRepo
}

// VehicleCompanyServiceDI ...
type VehicleCompanyServiceDI struct {
	dig.In
	Repo interfaces.VehicleCompanyRepo
}

// NewVehicleCompanyService ...
func NewVehicleCompanyService(di VehicleCompanyServiceDI) *VehicleCompanyService {
	return &VehicleCompanyService{
		repo: di.Repo,
	}
}

// FindByOgrnInn ...
func (s *VehicleCompanyService) FindByOgrnInn(ogrn int64, inn int64) (*entity.Company, error) {
	return s.repo.FindByOgrnInn(ogrn, inn)
}

// FindByINN ...
func (s *VehicleCompanyService) FindByINN(inn int64) (*entity.Company, error) {
	return s.repo.FindByINN(inn)
}

// FindByOGRN ...
func (s *VehicleCompanyService) FindByOGRN(ogrn int64) (*entity.Company, error) {
	return s.repo.FindByOGRN(ogrn)
}

// Create ...
func (s *VehicleCompanyService) Create(data *entity.Company) error {
	return s.repo.Create(data)
}

// Update ...
func (s *VehicleCompanyService) Update(data *entity.Company) error {
	return s.repo.Update(data)
}

// Delete ...
func (s *VehicleCompanyService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Read ...
func (s *VehicleCompanyService) Read(id int64) (*entity.Company, error) {
	return s.repo.Read(id)
}
