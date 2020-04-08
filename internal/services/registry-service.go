package services

import (
	"context"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"go.uber.org/dig"
	"sync"
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

// HandleParsed ...
func (s *RegistryService) HandleParsed(ctx context.Context, wg sync.WaitGroup, out chan interface{}) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case iface := <-out:
			value := iface.(string)
			fmt.Println(value)
		}
	}

}
