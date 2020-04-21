package services

import (
	"github.com/alexey-zayats/claim-parser/internal/config"
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"time"
)

// PeopleApplicationService ...
type PeopleApplicationService struct {
	config     *config.Config
	bidSvc     *PeopleBidService
	passSvc    *PeoplePassService
	companySvc *PeopleCompanyService
	branchSvc  *BranchService
}

// PeopleApplicationServiceDI ...
type PeopleApplicationServiceDI struct {
	dig.In
	Config     *config.Config
	BidSvc     *PeopleBidService
	PassSvc    *PeoplePassService
	CompanySvc *PeopleCompanyService
	BranchSvc  *BranchService
}

// NewPeopleApplicationService ...
func NewPeopleApplicationService(di PeopleApplicationServiceDI) *PeopleApplicationService {

	s := &PeopleApplicationService{
		config:     di.Config,
		bidSvc:     di.BidSvc,
		passSvc:    di.PassSvc,
		companySvc: di.CompanySvc,
		branchSvc:  di.BranchSvc,
	}

	return s
}

// SaveRecord ...
func (s *PeopleApplicationService) SaveRecord(a *model.Application) error {

	var err error
	var company *entity.CompanyPeople = nil

	company, err = s.companySvc.FindByINN(a.Inn)
	if err != nil {
		return errors.Wrapf(err, "unable find company by OGRN & INN")
	}

	if company == nil {
		company = &entity.CompanyPeople{
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

	bid := &entity.BidPeople{
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
	}

	if a.Dirty {
		bid.UserID = s.config.Pass.Dirty
	} else {
		bid.UserID = s.config.Pass.Clean
	}

	if err := s.bidSvc.Create(bid); err != nil {
		return errors.Wrap(err, "unable create bids_people record")
	}

	for _, pass := range a.Passes {

		p := &entity.PassPeople{
			BidID:      bid.ID,
			Source:     s.config.Pass.Source,
			DistrictID: a.DistrictID,
			PassType:   a.PassType,
			Lastname:   pass.Lastname,
			Firstname:  pass.Firstname,
			Patrname:   pass.Middlename,
			Shipping:   0,
		}

		if err := s.passSvc.Create(p); err != nil {
			return errors.Wrap(err, "unable create passes_people")
		}
	}

	return nil
}
