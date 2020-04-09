package services

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"time"
)

// IssuedService ...
type IssuedService struct {
	repo    interfaces.IssuedRepository
	passSvc *PassService
}

// IssuedServiceDI ...
type IssuedServiceDI struct {
	dig.In
	Repo    interfaces.IssuedRepository
	PassSvc *PassService
}

// NewIssuedService ...
func NewIssuedService(di IssuedServiceDI) *IssuedService {
	return &IssuedService{
		repo:    di.Repo,
		passSvc: di.PassSvc,
	}
}

// SaveRecord ...
func (s *IssuedService) SaveRecord(ctx context.Context, record *model.Registry) error {

	var issued *model.Issued
	var err error
	var create bool

	if issued, err = s.repo.FindByPass(record.PassNumber); err != nil {
		return errors.Wrap(err, "unable find issued")
	}

	if issued == nil {
		create = true
		issued = &model.Issued{}
	}

	issued.CreatedAt = time.Now()
	issued.CreatedBy = record.Event.CreatedBy
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

	pass, err := s.passSvc.FindByCar(issued.CompanyCar)
	if err != nil {
		return errors.Wrap(err, "unable find pass by issued.car")
	}

	if pass != nil {
		pass.IssuedID = issued.ID
		if err = s.passSvc.Update(pass); err != nil {
			return errors.Wrap(err, "unable Update pass status")
		}
	}

	return nil
}

// FindByCar ...
func (s *IssuedService) FindByCar(car string) (*model.Issued, error) {
	return s.repo.FindByCar(car)
}

// FindByPass ...
func (s *IssuedService) FindByPass(pass string) (*model.Issued, error) {
	return s.repo.FindByPass(pass)
}

// Create ...
func (s *IssuedService) Create(data *model.Issued) error {
	return s.repo.Create(data)
}

// Update ...
func (s *IssuedService) Update(data *model.Issued) error {
	return s.repo.Update(data)
}

// Delete ...
func (s *IssuedService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Read ...
func (s *IssuedService) Read(id int64) (*model.Issued, error) {
	return s.repo.Read(id)
}
