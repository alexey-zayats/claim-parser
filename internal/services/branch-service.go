package services

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// BranchService ...
type BranchService struct {
	repo interfaces.BranchRepository
}

// BranchServiceDI ...
type BranchServiceDI struct {
	dig.In
	Repo interfaces.BranchRepository
}

// NewBranchService ...
func NewBranchService(di BranchServiceDI) *BranchService {
	return &BranchService{
		repo: di.Repo,
	}
}

// Create ...
func (s *BranchService) Create(data *entity.Branch) error {
	return s.repo.Create(data)
}

// GetAll ...
func (s *BranchService) GetAll() (map[string]int64, error) {
	list, err := s.repo.GetAll()
	if err != nil {
		return nil, errors.Wrap(err, "unble get all branches")
	}

	table := make(map[string]int64)

	for _, branch := range list {
		table[branch.Name] = branch.ID
	}

	return table, nil
}
