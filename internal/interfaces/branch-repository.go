package interfaces

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
)

// BranchRepository ...
type BranchRepository interface {
	GetAll() ([]*entity.Branch, error)
	Create(data *entity.Branch) error
	FindByName(name string) (*entity.Branch, error)
}
