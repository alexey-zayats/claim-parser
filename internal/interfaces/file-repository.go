package interfaces

import "github.com/alexey-zayats/claim-parser/internal/model"

// FileRepository ...
type FileRepository interface {
	Update(*model.File) error
}
