package repository

import (
	"database/sql"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// RoutingRepository ...
type RoutingRepository struct {
	db *sqlx.DB
}

// RoutingRepositoryDI ...
type RoutingRepositoryDI struct {
	dig.In
	DB *sqlx.DB
}

// NewRoutingRepository ...
func NewRoutingRepository(di RoutingRepositoryDI) interfaces.RoutingRepository {
	return &RoutingRepository{
		db: di.DB,
	}
}

// FindBySourceDistrict ...
func (r *RoutingRepository) FindBySourceDistrict(sourceID, districtID int64) (*entity.Routing, error) {
	var record entity.Routing

	query := "SELECT id, source_id, district_id, clean_id, dirty_id " +
		"FROM routing WHERE source_id = ? AND district_id = ?"

	err := r.db.Get(&record, query, sourceID, districtID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "unable find routing record by source(%d) and district(%d)",
			sourceID, districtID)
	}

	return &record, nil
}
