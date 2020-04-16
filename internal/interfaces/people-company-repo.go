package interfaces

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
)

// PeopleCompanyRepo ...
type PeopleCompanyRepo interface {
	Create(people *entity.CompanyPeople) error
	Read(id int64) (*entity.CompanyPeople, error)
	Update(people *entity.CompanyPeople) error
	Delete(id int64) error

	FindByINN(inn int64) (*entity.CompanyPeople, error)
	FindByOGRN(inn int64) (*entity.CompanyPeople, error)
	FindByOgrnInn(ogrn, inn int64) (*entity.CompanyPeople, error)
}
