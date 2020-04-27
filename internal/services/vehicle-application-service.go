package services

import (
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/application"
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"time"
)

// VehicleApplicationService ...
type VehicleApplicationService struct {
	config     *config.Config
	bidSvc     *VehicleBidService
	passSvc    *VehiclePassService
	companySvc *VehicleCompanyService
	branchSvc  *BranchService
	sourceSvc  *SourceService
	routingSvc *RoutingService
}

// VehicleApplicationServiceDI ...
type VehicleApplicationServiceDI struct {
	dig.In
	Config     *config.Config
	BidSvc     *VehicleBidService
	PassSvc    *VehiclePassService
	CompanySvc *VehicleCompanyService
	BranchSvc  *BranchService
	SourceSvc  *SourceService
	RoutingSvc *RoutingService
}

// NewVehicleApplicationService ...
func NewVehicleApplicationService(di VehicleApplicationServiceDI) *VehicleApplicationService {

	s := &VehicleApplicationService{
		config:     di.Config,
		bidSvc:     di.BidSvc,
		passSvc:    di.PassSvc,
		companySvc: di.CompanySvc,
		branchSvc:  di.BranchSvc,
		sourceSvc:  di.SourceSvc,
		routingSvc: di.RoutingSvc,
	}

	return s
}

// SaveRecord ...
func (s *VehicleApplicationService) SaveRecord(a *application.Vehicle) error {

	var err error
	var company *entity.Company = nil

	company, err = s.companySvc.FindByINN(a.Inn)
	if err != nil {
		return errors.Wrapf(err, "unable find company by OGRN & INN")
	}

	if company == nil {
		company = &entity.Company{
			OGRN:     a.Ogrn,
			INN:      a.Inn,
			Name:     a.Title,
			BranchID: a.ActivityKind,
			Status:   0,
		}
		if err = s.companySvc.Create(company); err != nil {
			return errors.Wrapf(err, "unable create company")
		}
	} else {

		update := false
		if company.INN == 0 {
			company.INN = a.Inn
			update = true
		}

		if company.OGRN == 0 {
			company.OGRN = a.Ogrn
			update = true
		}

		if update {
			if err = s.companySvc.Update(company); err != nil {
				return errors.Wrapf(err, "unable update company")
			}
		}
	}

	sourceName := "form.vehicle"
	source, err := s.sourceSvc.FindByName(sourceName)
	if err != nil {
		return errors.Wrapf(err, "unable to find source with name %s", sourceName)
	}
	if source == nil {
		return fmt.Errorf("unable to find source with name %s", sourceName)
	}

	routing, err := s.routingSvc.FindBySourceDistrict(source.ID, a.DistrictID)
	if err != nil {
		return errors.Wrapf(err, "unable to find routing by source(%d) and district(%d)", source.ID, a.DistrictID)
	}
	if routing == nil {
		return fmt.Errorf("unable to find routing by source(%d) and district(%d)", source.ID, a.DistrictID)
	}

	userID := int64(0)

	if a.Dirty {
		userID = routing.DirtyID
	} else {
		userID = routing.CleanID
	}

	t := time.Now()

	bid := &entity.Bid{
		CompanyID:       company.ID,
		BranchID:        a.ActivityKind,
		CompanyName:     a.Title,
		CompanyAddress:  a.Address,
		CompanyCeoPhone: a.CeoPhone,
		CompanyCeoEmail: a.CeoEmail,
		CompanyCeoName:  a.CeoName,
		Agree:           a.Agreement,
		Confirm:         a.Reliability,
		WorkflowStatus:  1,
		DistrictID:      a.DistrictID,
		PassType:        a.PassType,
		CreatedAt:       time.Now(),
		CreatedBy:       userID,
		DateFrom:        t,
		DateTo:          t,
	}

	if err := s.bidSvc.Create(bid); err != nil {
		return errors.Wrap(err, "unable create bids record")
	}

	for _, pass := range a.Passes {

		p := &entity.Pass{
			BidID:      bid.ID,
			Source:     source.ID,
			DistrictID: a.DistrictID,
			PassType:   a.PassType,
			Car:        pass.Car,
			Lastname:   pass.Lastname,
			Firstname:  pass.Firstname,
			Patrname:   pass.Middlename,
			Shipping:   0,
			CreatedAt:  time.Now(),
			CreatedBy:  userID,
		}

		if err := s.passSvc.Create(p); err != nil {
			return errors.Wrap(err, "unable create passes")
		}
	}

	return nil
}
