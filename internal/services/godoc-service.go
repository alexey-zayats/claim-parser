package services

import (
	"context"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// GodocService ...
type GodocService struct {
	bidSvc  *BidService
	passSvc *PassService
}

// GodocServiceDI ...
type GodocServiceDI struct {
	dig.In
	BidSvc  *BidService
	PassSvc *PassService
}

// NewGodocService ...
func NewGodocService(input GodocServiceDI) *GodocService {
	return &GodocService{
		bidSvc:  input.BidSvc,
		passSvc: input.PassSvc,
	}
}

// SaveClaim ...
func (s *GodocService) SaveClaim(ctx context.Context, claim *model.Claim) error {

	bid := &model.Bid{
		WorkflowStatus: 1,
		Code:           claim.Code,
		CreatedAt:      claim.Created,
		District:       claim.DistrictID,
		Source:         claim.Source,
	}

	if err := s.bidSvc.Create(bid); err != nil {
		return errors.Wrap(err, "unable create bid")
	}

	for _, car := range claim.Cars {

		select {
		case <-ctx.Done():
			return fmt.Errorf("canceled")
		default:

			pass := &model.Pass{
				CompanyBranch:     claim.Company.Activity,
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
				Source:            3,
				District:          claim.DistrictID,
				Status:            0,
				CreatedAt:         claim.Created,
				CreatedBy:         1,
				BidID:             bid.ID,
			}

			if err := s.passSvc.Create(pass); err != nil {
				return errors.Wrap(err, "unable create pass")
			}
		}
	}

	return nil
}
