package services

import (
	"github.com/alexey-zayats/claim-parser/internal/entity"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
)

// VehicleClaimService ...
type VehicleClaimService struct {
	bidSvc     *VehicleBidService
	passSvc    *VehiclePassService
	issuedSvc  *VehicleIssuedService
	companySvc *VehicleCompanyService
	branchSvc  *BranchService
	branches   map[string]int64
}

// VehicleClaimServiceDI ...
type VehicleClaimServiceDI struct {
	dig.In
	BidSvc     *VehicleBidService
	PassSvc    *VehiclePassService
	IssuedSvc  *VehicleIssuedService
	CompanySvc *VehicleCompanyService
	BranchSvc  *BranchService
}

// NewVehicleClaimService ...
func NewVehicleClaimService(di VehicleClaimServiceDI) *VehicleClaimService {
	s := &VehicleClaimService{
		bidSvc:     di.BidSvc,
		passSvc:    di.PassSvc,
		issuedSvc:  di.IssuedSvc,
		companySvc: di.CompanySvc,
		branchSvc:  di.BranchSvc,
	}

	var err error
	s.branches, err = s.branchSvc.GetAll()
	if err != nil {
		logrus.WithFields(logrus.Fields{"reason": err}).Error("unable get branches")
	}

	return s
}

// SaveRecord ...
func (s *VehicleClaimService) SaveRecord(event *model.Event, claim *model.VehicleClaim) error {

	var err error
	var company *entity.Company = nil

	//company, err = s.companySvc.FindByOgrnInn(claim.Company.PSRN, claim.Company.TIN)
	//if err != nil {
	//	return errors.Wrapf(err, "unable find company by OGRN & INN")
	//}
	//
	//if company == nil {
	company, err = s.companySvc.FindByINN(claim.Company.TIN)
	if err != nil {
		return errors.Wrapf(err, "unable find company by OGRN & INN")
	}
	//}
	//
	//if company == nil {
	//	company, err = s.companySvc.FindByOGRN(claim.Company.PSRN)
	//	if err != nil {
	//		return errors.Wrapf(err, "unable find company by OGRN & INN")
	//	}
	//}

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
	} else {

		if company.INN == 0 {
			company.INN = claim.Company.TIN
		}

		if company.OGRN == 0 {
			company.OGRN = claim.Company.PSRN
		}

		if err = s.companySvc.Update(company); err != nil {
			return errors.Wrapf(err, "unable create company")
		}
	}

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
		UserID:          event.CreatedBy,
		Source:          claim.Source,
		Agree:           1,
		Confirm:         1,
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
			Source:     event.Source,
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
