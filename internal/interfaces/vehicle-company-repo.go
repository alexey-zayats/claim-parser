package interfaces

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
)

// VehicleCompanyRepo ...
type VehicleCompanyRepo interface {
	Create(people *entity.Company) error
	Read(id int64) (*entity.Company, error)
	Update(people *entity.Company) error
	Delete(id int64) error

	FindByINN(inn string) (*entity.Company, error)
	FindByOGRN(inn string) (*entity.Company, error)
	FindByOgrnInn(ogrn, inn string) (*entity.Company, error)
}
