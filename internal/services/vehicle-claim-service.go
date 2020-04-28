package services

import (
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"time"
)

// VehicleClaimService ...
type VehicleClaimService struct {
	bidSvc     *VehicleBidService
	passSvc    *VehiclePassService
	issuedSvc  *VehicleIssuedService
	companySvc *VehicleCompanyService
	branchSvc  *BranchService
	sourceSvc  *SourceService
	routingSvc *RoutingService
}

// VehicleClaimServiceDI ...
type VehicleClaimServiceDI struct {
	dig.In
	BidSvc     *VehicleBidService
	PassSvc    *VehiclePassService
	IssuedSvc  *VehicleIssuedService
	CompanySvc *VehicleCompanyService
	BranchSvc  *BranchService
	SourceSvc  *SourceService
	RoutingSvc *RoutingService
}

// NewVehicleClaimService ...
func NewVehicleClaimService(di VehicleClaimServiceDI) *VehicleClaimService {
	s := &VehicleClaimService{
		bidSvc:     di.BidSvc,
		passSvc:    di.PassSvc,
		issuedSvc:  di.IssuedSvc,
		companySvc: di.CompanySvc,
		branchSvc:  di.BranchSvc,
		sourceSvc:  di.SourceSvc,
		routingSvc: di.RoutingSvc,
	}

	return s
}

// SaveRecord ...
func (s *VehicleClaimService) SaveRecord(event *model.Event, claim *model.VehicleClaim) error {

	var err error
	var company *entity.Company = nil

	company, err = s.companySvc.FindByINN(claim.Company.TIN)
	if err != nil {
		return errors.Wrapf(err, "unable find company by OGRN & INN")
	}

	branchID := event.BranchID

	if branchID == 0 {
		branch, err := s.branchSvc.FindByName(claim.Company.Activity)
		if err != nil {
			return errors.Wrapf(err, "unable find branch by name")
		}
		if branch == nil {
			branch = &entity.Branch{
				Name: claim.Company.Activity,
				Type: "Произвольные",
			}
			if err := s.branchSvc.Create(branch); err != nil {
				return errors.Wrapf(err, "unable create branch")
			}
		}
		branchID = branch.ID
	}

	if company == nil {
		company = &entity.Company{
			OGRN:     claim.Company.PSRN,
			INN:      claim.Company.TIN,
			Name:     claim.Company.Title,
			BranchID: branchID,
			Status:   0,
		}
		if err = s.companySvc.Create(company); err != nil {
			return errors.Wrapf(err, "unable create company")
		}
	} else if company.OGRN == 0 {

		company.OGRN = claim.Company.PSRN

		if err = s.companySvc.Update(company); err != nil {
			return errors.Wrapf(err, "unable create company")
		}
	} else {
		logrus.WithFields(logrus.Fields{
			"inn":  company.INN,
			"ogrn": company.OGRN,
		}).Debug("company")
	}

	sourceName := "gsheet-vehicle"
	source, err := s.sourceSvc.FindByName(sourceName)
	if err != nil {
		return errors.Wrapf(err, "unable to find source with name %s", sourceName)
	}
	if source == nil {
		return fmt.Errorf("unable to find source with name %s", sourceName)
	}

	routing, err := s.routingSvc.FindBySourceDistrict(source.ID, event.DistrictID)
	if err != nil {
		return errors.Wrapf(err, "unable to find routing by source(%d) and district(%d)", source.ID, event.DistrictID)
	}
	if routing == nil {
		return fmt.Errorf("unable to find routing by source(%d) and district(%d)", source.ID, event.DistrictID)
	}

	userID := int64(0)

	if claim.Success {
		userID = routing.DirtyID
	} else {
		userID = routing.CleanID
	}

	t := time.Now()

	bid := &entity.Bid{
		CompanyID:       company.ID,
		FileID:          event.FileID,
		WorkflowStatus:  1,
		Code:            claim.Code,
		BranchID:        branchID,
		DistrictID:      event.DistrictID,
		CompanyBranch:   claim.Company.Activity,
		CompanyName:     claim.Company.Title,
		CompanyAddress:  claim.Company.Address,
		CompanyCeoPhone: claim.Company.HeadPhone,
		CompanyCeoEmail: claim.Company.HeadEmail,
		CompanyCeoName:  claim.Company.HeadName,
		PassType:        event.PassType,
		CreatedAt:       claim.Created,
		CreatedBy:       event.CreatedBy,
		UserID:          userID,
		Source:          claim.Source,
		Agree:           1,
		Confirm:         1,
		DateFrom:        t,
		DateTo:          t,
	}

	if err := s.bidSvc.Create(bid); err != nil {
		return errors.Wrap(err, "unable create bids record")
	}

	for _, pass := range claim.Passes {

		p := &entity.Pass{
			BidID:      bid.ID,
			Lastname:   pass.FIO.Lastname,
			Firstname:  pass.FIO.Firstname,
			Patrname:   pass.FIO.Patronymic,
			Car:        pass.Number,
			Source:     source.ID,
			DistrictID: event.DistrictID,
			PassType:   event.PassType,
			Status:     0,
			FileID:     event.FileID,
			CreatedAt:  claim.Created,
			CreatedBy:  event.CreatedBy,
		}

		if event.Check == 1 {
			issued, err := s.issuedSvc.FindByCar(pass.Number)
			if err != nil {
				return errors.Wrap(err, "unable find issued by car")
			}
			if issued != nil {
				// FIXME: нужен правильный ID статуса полученного пропуска
				p.Status = 100
				p.IssuedID = issued.ID
			}
		}

		if err := s.passSvc.Create(p); err != nil {
			return errors.Wrap(err, "unable create pass")
		}
	}

	return nil
}
