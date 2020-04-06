package interfaces

import "github.com/alexey-zayats/claim-parser/internal/model"

// PassRepository ...
type PassRepository interface {
	Create(*model.Pass) error
}
