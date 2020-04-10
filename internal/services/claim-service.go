package services

import (
	"context"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// ClaimService ...
type ClaimService struct {
	bidSvc    *BidService
	passSvc   *PassService
	issuedSvc *IssuedService
}

// ClaimServiceDI ...
type ClaimServiceDI struct {
	dig.In
	BidSvc    *BidService
	PassSvc   *PassService
	IssuedSvc *IssuedService
}

// NewClaimService ...
func NewClaimService(di ClaimServiceDI) *ClaimService {
	return &ClaimService{
		bidSvc:    di.BidSvc,
		passSvc:   di.PassSvc,
		issuedSvc: di.IssuedSvc,
	}
}

// SaveRecord ...
func (s *ClaimService) SaveRecord(ctx context.Context, claim *model.Claim) error {

	event := claim.Event

	bid := &model.Bid{
		FileID:         event.FileID,
		WorkflowStatus: 1,
		Code:           claim.Code,
		District:       event.District,
		PassType:       event.PassType,
		CreatedAt:      claim.Created,
		CreatedBy:      event.CreatedBy,
		UserID:         event.CreatedBy,
		Source:         claim.Source,
	}

	if err := s.bidSvc.Create(bid); err != nil {
		return errors.Wrap(err, "unable create bids record")
	}

	for _, car := range claim.Cars {

		p := &model.Pass{
			CompanyBranch:     claim.Company.Activity,
			CompanyOkved:      "",
			CompanyInn:        claim.Company.INN,
			CompanyName:       claim.Company.Title,
			CompanyAddress:    claim.Company.Address,
			CompanyCeoPhone:   claim.Company.Head.Contact.Phone,
			CompanyCeoEmail:   claim.Company.Head.Contact.EMail,
			CompanyLastname:   claim.Company.Head.FIO.Surname,
			CompanyFirstname:  claim.Company.Head.FIO.Name,
			CompanyPatrname:   claim.Company.Head.FIO.Patronymic,
			EmployeeLastname:  car.FIO.Surname,
			EmployeeFirstname: car.FIO.Name,
			EmployeePatrname:  car.FIO.Patronymic,
			EmployeeCar:       car.Number,
			EmployeeAgree:     1,
			EmployeeConfirm:   1,
			Source:            event.Source,
			District:          event.District,
			PassType:          event.PassType,
			Status:            0,
			FileID:            event.FileID,
			CreatedAt:         claim.Created,
			CreatedBy:         event.CreatedBy,
			BidID:             bid.ID,
			Ogrn:              claim.Ogrn,
		}

		if event.Check == 1 {
			issued, err := s.issuedSvc.FindByCar(car.Number)
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
