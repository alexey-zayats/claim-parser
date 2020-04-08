package services

import (
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"go.uber.org/dig"
)

// RegistryService ...
type RegistryService struct {
	passRepo interfaces.PassRepository
}

// RegistryServiceDI ...
type RegistryServiceDI struct {
	dig.In
	PassRepo interfaces.PassRepository
}

// NewRegistryService ...
func NewRegistryService(input RegistryServiceDI) *RegistryService {
	return &RegistryService{
		passRepo: input.PassRepo,
	}
}
