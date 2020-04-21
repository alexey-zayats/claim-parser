package services

import (
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/model"
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
}

// VehicleApplicationServiceDI ...
type VehicleApplicationServiceDI struct {
	dig.In
	Config     *config.Config
	BidSvc     *VehicleBidService
	PassSvc    *VehiclePassService
	CompanySvc *VehicleCompanyService
	BranchSvc  *BranchService
}

// NewVehicleApplicationService ...
func NewVehicleApplicationService(di VehicleApplicationServiceDI) *VehicleApplicationService {

	s := &VehicleApplicationService{
		config:     di.Config,
		bidSvc:     di.BidSvc,
		passSvc:    di.PassSvc,
		companySvc: di.CompanySvc,
		branchSvc:  di.BranchSvc,
	}

	return s
}

// SaveRecord ...
func (s *VehicleApplicationService) SaveRecord(a *model.Application) error {

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

	userID := int64(0)

	if a.Dirty {
		userID = s.config.Pass.Dirty
	} else {
		userID = s.config.Pass.Clean
	}

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
		CreatedBy:       s.config.Pass.Creator,
		UserID:          userID,
	}

	if err := s.bidSvc.Create(bid); err != nil {
		return errors.Wrap(err, "unable create bids record")
	}

	for _, pass := range a.Passes {

		p := &entity.Pass{
			BidID:      bid.ID,
			Source:     s.config.Pass.Source,
			DistrictID: a.DistrictID,
			PassType:   a.PassType,
			Car:        pass.Car,
			Lastname:   pass.Lastname,
			Firstname:  pass.Firstname,
			Patrname:   pass.Middlename,
			Shipping:   0,
			CreatedAt:  time.Now(),
			CreatedBy:  s.config.Pass.Creator,
		}

		if err := s.passSvc.Create(p); err != nil {
			return errors.Wrap(err, "unable create passes")
		}
	}

	return nil
}
