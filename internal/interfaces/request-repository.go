package interfaces

import "github.com/alexey-zayats/claim-parser/internal/model"

// RequestRepository ...
type RequestRepository interface {
	Create(*model.Request) (int64, error)
	Read(id int) (*model.Request, error)
	Update(*model.Request) error
	Delete(id int) error
}
