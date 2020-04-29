package services

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"go.uber.org/dig"
)

// PeopleCompanyService ...
type PeopleCompanyService struct {
	repo interfaces.PeopleCompanyRepo
}

// PeopleCompanyServiceDI ...
type PeopleCompanyServiceDI struct {
	dig.In
	Repo interfaces.PeopleCompanyRepo
}

// NewPeopleCompanyService ...
func NewPeopleCompanyService(di PeopleCompanyServiceDI) *PeopleCompanyService {
	return &PeopleCompanyService{
		repo: di.Repo,
	}
}

// FindByOgrnInn ...
func (s *PeopleCompanyService) FindByOgrnInn(ogrn, inn string) (*entity.CompanyPeople, error) {
	return s.repo.FindByOgrnInn(ogrn, inn)
}

// FindByINN ...
func (s *PeopleCompanyService) FindByINN(inn string) (*entity.CompanyPeople, error) {
	return s.repo.FindByINN(inn)
}

// FindByOGRN ...
func (s *PeopleCompanyService) FindByOGRN(ogrn string) (*entity.CompanyPeople, error) {
	return s.repo.FindByOGRN(ogrn)
}

// Create ...
func (s *PeopleCompanyService) Create(data *entity.CompanyPeople) error {
	return s.repo.Create(data)
}

// Update ...
func (s *PeopleCompanyService) Update(data *entity.CompanyPeople) error {
	return s.repo.Update(data)
}

// Delete ...
func (s *PeopleCompanyService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Read ...
func (s *PeopleCompanyService) Read(id int64) (*entity.CompanyPeople, error) {
	return s.repo.Read(id)
}
