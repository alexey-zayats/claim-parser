package services

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"go.uber.org/dig"
)

// PeoplePassService ...
type PeoplePassService struct {
	repo interfaces.PeoplePassRepo
}

// PeoplePassServiceDI ...
type PeoplePassServiceDI struct {
	dig.In
	Repo interfaces.PeoplePassRepo
}

// NewPeoplePassService ...
func NewPeoplePassService(di PeoplePassServiceDI) *PeoplePassService {
	return &PeoplePassService{
		repo: di.Repo,
	}
}

// Create ...
func (s *PeoplePassService) Create(data *entity.PassPeople) error {
	return s.repo.Create(data)
}

// Update ...
func (s *PeoplePassService) Update(data *entity.PassPeople) error {
	return s.repo.Update(data)
}

// Delete ...
func (s *PeoplePassService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Read ...
func (s *PeoplePassService) Read(id int64) (*entity.PassPeople, error) {
	return s.repo.Read(id)
}
