package services

import (
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// PeopleClaimService ...
type PeopleClaimService struct {
	bidSvc     *PeopleBidService
	passSvc    *PeoplePassService
	companySvc *PeopleCompanyService
	branchSvc  *BranchService
	sourceSvc  *SourceService
	routingSvc *RoutingService
}

// PeopleClaimServiceDI ...
type PeopleClaimServiceDI struct {
	dig.In
	BidSvc     *PeopleBidService
	PassSvc    *PeoplePassService
	CompanySvc *PeopleCompanyService
	BranchSvc  *BranchService
	SourceSvc  *SourceService
	RoutingSvc *RoutingService
}

// NewPeopleClaimService ...
func NewPeopleClaimService(di PeopleClaimServiceDI) *PeopleClaimService {

	s := &PeopleClaimService{
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
func (s *PeopleClaimService) SaveRecord(event *model.Event, claim *model.PeopleClaim) error {

	var err error
	var company *entity.CompanyPeople = nil

	company, err = s.companySvc.FindByINN(claim.Company.INN)
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
			//
		}
		branchID = branch.ID
	}

	if company == nil {
		company = &entity.CompanyPeople{
			OGRN:     claim.Company.OGRN,
			INN:      claim.Company.INN,
			Name:     claim.Company.Title,
			BranchID: branchID,
			Status:   0,
		}
		if err = s.companySvc.Create(company); err != nil {
			return errors.Wrapf(err, "unable create company")
		}
	} else if len(company.OGRN) == 0 {

		company.OGRN = claim.Company.OGRN

		if err = s.companySvc.Update(company); err != nil {
			return errors.Wrapf(err, "unable create company")
		}
	}

	sourceName := "gsheet-people"
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

	bid := &entity.BidPeople{
		FileID:          event.FileID,
		CompanyID:       company.ID,
		BranchID:        branchID,
		CompanyBranch:   claim.Company.Activity,
		CompanyName:     claim.Company.Title,
		CompanyAddress:  claim.Company.Address,
		CompanyCeoPhone: claim.Company.HeadPhone,
		CompanyCeoEmail: claim.Company.HeadEmail,
		CompanyCeoName:  claim.Company.HeadName,
		Agree:           1,
		Confirm:         1,
		WorkflowStatus:  1,
		DistrictID:      event.DistrictID,
		PassType:        event.PassType,
		Source:          claim.Source,
		CreatedAt:       claim.Created,
		CreatedBy:       event.CreatedBy,
		MovedTo:         userID,
	}

	if err := s.bidSvc.Create(bid); err != nil {
		return errors.Wrap(err, "unable create bids_people record")
	}

	for _, pass := range claim.Passes {

		p := &entity.PassPeople{
			BidID:      bid.ID,
			Source:     source.ID,
			DistrictID: event.DistrictID,
			PassType:   event.PassType,
			Lastname:   pass.FIO.Lastname,
			Firstname:  pass.FIO.Firstname,
			Patrname:   pass.FIO.Patronymic,
			Shipping:   0,
		}

		if err := s.passSvc.Create(p); err != nil {
			return errors.Wrap(err, "unable create passes_people")
		}
	}

	return nil
}
