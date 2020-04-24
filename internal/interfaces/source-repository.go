package interfaces

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
)

// SourceRepository ...
type SourceRepository interface {
	FindByName(name string) (*entity.Source, error)
}
