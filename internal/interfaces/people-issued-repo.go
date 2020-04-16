package interfaces

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
)

// PeopleIssuedRepo ...
type PeopleIssuedRepo interface {
	Create(data *entity.IssuedPeople) error
	Read(id int64) (*entity.IssuedPeople, error)
	Update(data *entity.IssuedPeople) error
	Delete(id int64) error
}
