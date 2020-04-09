package services

import (
	"context"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// GodocService ...
type GodocService struct {
	passRepo interfaces.PassRepository
	reqRepo  interfaces.BidRepository
}

// GodocServiceDI ...
type GodocServiceDI struct {
	dig.In
	PassRepo interfaces.PassRepository
	ReqRepo  interfaces.BidRepository
}

// NewGodocService ...
func NewGodocService(input GodocServiceDI) *GodocService {
	return &GodocService{
		passRepo: input.PassRepo,
		reqRepo:  input.ReqRepo,
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

	id, err := s.reqRepo.Create(bid)
	if err != nil {
		return errors.Wrap(err, "unable create bid")
	}
	bid.ID = int(id)

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

			id, err = s.passRepo.Create(pass)
			if err != nil {
				return errors.Wrap(err, "unable create pass")
			}
			pass.ID = int(id)
		}
	}

	return nil
}
