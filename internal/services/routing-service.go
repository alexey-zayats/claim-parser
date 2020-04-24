package services

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"go.uber.org/dig"
)

// RoutingService ...
type RoutingService struct {
	repo interfaces.RoutingRepository
}

// RoutingServiceDI ...
type RoutingServiceDI struct {
	dig.In
	Repo interfaces.RoutingRepository
}

// NewRoutingService ...
func NewRoutingService(di RoutingServiceDI) *RoutingService {
	return &RoutingService{
		repo: di.Repo,
	}
}

// FindBySourceDistrict ...
func (s *RoutingService) FindBySourceDistrict(sourceID, districtID int64) (*entity.Routing, error) {
	return s.repo.FindBySourceDistrict(sourceID, districtID)
}
