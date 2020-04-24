package services

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"go.uber.org/dig"
)

// SourceService ...
type SourceService struct {
	repo interfaces.SourceRepository
}

// SourceServiceDI ...
type SourceServiceDI struct {
	dig.In
	Repo interfaces.SourceRepository
}

// NewSourceService ...
func NewSourceService(di SourceServiceDI) *SourceService {
	return &SourceService{
		repo: di.Repo,
	}
}

// FindByName ...
func (s *SourceService) FindByName(name string) (*entity.Source, error) {
	return s.repo.FindByName(name)
}
