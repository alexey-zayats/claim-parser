package services

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"time"
)

// VehicleIssuedService ...
type VehicleIssuedService struct {
	repo    interfaces.VehicleIssuedRepo
	passSvc *VehiclePassService
}

// VehicleIssuedServiceDI ...
type VehicleIssuedServiceDI struct {
	dig.In
	Repo    interfaces.VehicleIssuedRepo
	PassSvc *VehiclePassService
}

// NewVehicleIssuedService ...
func NewVehicleIssuedService(di VehicleIssuedServiceDI) *VehicleIssuedService {
	return &VehicleIssuedService{
		repo:    di.Repo,
		passSvc: di.PassSvc,
	}
}

// SaveRecord ...
func (s *VehicleIssuedService) SaveRecord(event *model.Event, record *model.VehicleRegistry) error {

	var issued *entity.Issued
	var err error
	var create bool

	if issued, err = s.repo.FindByPass(record.PassNumber); err != nil {
		return errors.Wrap(err, "unable find issued")
	}

	if issued == nil {
		create = true
		issued = &entity.Issued{}
	}

	issued.FileID = event.FileID
	issued.CreatedAt = time.Now()
	issued.CreatedBy = event.CreatedBy
	issued.CompanyInn = record.CompanyInn
	issued.CompanyOgrn = record.CompanyOgrn
	issued.CompanyName = record.CompanyName
	issued.CompanyFio = record.CompanyFio
	issued.CompanyCar = record.CompanyCar
	issued.LegalBasement = record.LegalBasement
	issued.PassNumber = record.PassNumber
	issued.District = record.District
	issued.PassType = record.PassType
	issued.IssuedAt = record.IssuedAt
	issued.RegistryNumber = record.RegistryNumber
	issued.Shipping = record.Shipping

	if create {
		if err := s.repo.Create(issued); err != nil {

			logrus.WithFields(logrus.Fields{
				"company": issued.CompanyName,
				"car":     issued.CompanyCar,
				"pass":    issued.PassNumber}).Error("registry")

			return errors.Wrap(err, "unable create issued record")
		}
	} else {
		if err := s.repo.Update(issued); err != nil {

			logrus.WithFields(logrus.Fields{
				"company": issued.CompanyName,
				"car":     issued.CompanyCar,
				"pass":    issued.PassNumber}).Error("registry")

			return errors.Wrap(err, "unable Update issued record")
		}
	}

	if len(issued.CompanyCar) > 0 {

		pass, err := s.passSvc.FindByCar(issued.CompanyCar)
		if err != nil {
			return errors.Wrap(err, "unable find pass by issued.car")
		}

		if pass != nil {

			// FIXME: нужен правильный ID статуса полученного пропуска
			pass.Status = 100
			pass.IssuedID = issued.ID

			if err = s.passSvc.Update(pass); err != nil {
				return errors.Wrap(err, "unable Update pass status")
			}
		}

	}

	return nil
}

// FindByCar ...
func (s *VehicleIssuedService) FindByCar(car string) (*entity.Issued, error) {
	return s.repo.FindByCar(car)
}

// FindByPass ...
func (s *VehicleIssuedService) FindByPass(pass string) (*entity.Issued, error) {
	return s.repo.FindByPass(pass)
}

// Create ...
func (s *VehicleIssuedService) Create(data *entity.Issued) error {
	return s.repo.Create(data)
}

// Update ...
func (s *VehicleIssuedService) Update(data *entity.Issued) error {
	return s.repo.Update(data)
}

// Delete ...
func (s *VehicleIssuedService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Read ...
func (s *VehicleIssuedService) Read(id int64) (*entity.Issued, error) {
	return s.repo.Read(id)
}
