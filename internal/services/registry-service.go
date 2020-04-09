package services

import (
	"go.uber.org/dig"
)

// RegistryService ...
type RegistryService struct {
	passSvc *PassService
	bidSvc  *BidService
}

// RegistryServiceDI ...
type RegistryServiceDI struct {
	dig.In
	PassSvc *PassService
	BidSvc  *BidService
}

// NewRegistryService ...
func NewRegistryService(di RegistryServiceDI) *RegistryService {
	return &RegistryService{
		passSvc: di.PassSvc,
		bidSvc:  di.BidSvc,
	}
}
