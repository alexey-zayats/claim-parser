package interfaces

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
)

// PeoplePassRepo ...
type PeoplePassRepo interface {
	Create(people *entity.PassPeople) error
	Read(id int64) (*entity.PassPeople, error)
	Update(*entity.PassPeople) error
	Delete(id int64) error
}
