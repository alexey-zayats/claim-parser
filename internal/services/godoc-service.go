package services

import (
	"context"
	"fmt"
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"sync"
)

// GodocService ...
type GodocService struct {
	passRepo interfaces.PassRepository
	reqRepo  interfaces.RequestRepository
}

// GodocServiceDI ...
type GodocServiceDI struct {
	dig.In
	PassRepo interfaces.PassRepository
	ReqRepo  interfaces.RequestRepository
}

// NewGodocService ...
func NewGodocService(input GodocServiceDI) *GodocService {
	return &GodocService{
		passRepo: input.PassRepo,
		reqRepo:  input.ReqRepo,
	}
}

// HandleParsed ...
func (s *GodocService) HandleParsed(ctx context.Context, wg sync.WaitGroup, out chan interface{}) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case iface := <-out:

			claim := iface.(*model.Claim)

			logrus.WithFields(logrus.Fields{"company": claim.Company.Title}).Debug("claim")

			if err := s.SaveClaim(ctx, claim); err != nil {
				logrus.WithFields(logrus.Fields{"reason": err}).Error("unable save claim")
			}
		}
	}
}

// SaveClaim ...
func (s *GodocService) SaveClaim(ctx context.Context, claim *model.Claim) error {

	req := &model.Request{
		Status:         0,
		WorkflowStatus: 1,
		Code:           claim.Code,
		CreatedAt:      claim.Created,
		District:       claim.DistrictID,
		Source:         claim.Source,
	}

	id, err := s.reqRepo.Create(req)
	if err != nil {
		return errors.Wrap(err, "unable create bid")
	}
	req.ID = int(id)

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
				SendType:          "formstruct-dump",
				Status:            0,
				CreatedAt:         claim.Created,
				CreatedBy:         1,
				RequestID:         req.ID,
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
