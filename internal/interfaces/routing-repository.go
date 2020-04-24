package interfaces

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
)

// RoutingRepository ...
type RoutingRepository interface {
	FindBySourceDistrict(sourceID, districtID int64) (*entity.Routing, error)
}
