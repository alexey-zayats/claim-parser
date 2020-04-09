package services

import (
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"go.uber.org/dig"
)

// PassService ...
type PassService struct {
	repo interfaces.PassRepository
}

// PassServiceDI ...
type PassServiceDI struct {
	dig.In
	Repo interfaces.PassRepository
}

// NewPassService ...
func NewPassService(di PassServiceDI) *PassService {
	return &PassService{
		repo: di.Repo,
	}
}

// Create ...
func (s *PassService) Create(data *model.Pass) error {
	return s.repo.Create(data)
}

// Update ...
func (s *PassService) Update(data *model.Pass) error {
	return s.repo.Update(data)
}

// Delete ...
func (s *PassService) Delete(id int64) error {
	return s.repo.Delete(id)
}

// Read ...
func (s *PassService) Read(id int64) (*model.Pass, error) {
	return s.repo.Read(id)
}

// FindByCar ...
func (s *PassService) FindByCar(car string) (*model.Pass, error) {
	return s.repo.FindByCar(car)
}
