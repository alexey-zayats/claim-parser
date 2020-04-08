package services

import (
	"github.com/alexey-zayats/claim-parser/internal/interfaces"
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"time"
)

// EventService ...
type EventService struct {
	fileRepo interfaces.FileRepository
	passRepo interfaces.PassRepository
	reqRepo  interfaces.RequestRepository
}

// EventServiceInput ...
type EventServiceInput struct {
	dig.In
	FileRepo interfaces.FileRepository
	PassRepo interfaces.PassRepository
	ReqRepo  interfaces.RequestRepository
}

// NewEventService ...
func NewEventService(input EventServiceInput) *EventService {
	return &EventService{
		fileRepo: input.FileRepo,
		passRepo: input.PassRepo,
		reqRepo:  input.ReqRepo,
	}
}

// StoreEvent ...
func (s *EventService) StoreEvent(e *model.Event) {

	r := &model.Request{
		FileID:         e.FileID,
		Status:         0,
		WorkflowStatus: 1,
		Code:           e.Claim.Code,
		District:       e.District,
		PassType:       e.PassType,
		CreatedAt:      time.Now(),
		CreatedBy:      e.CreatedBy,
		UserID:         e.CreatedBy,
		Source:         e.Claim.Source,
	}

	id, err := s.reqRepo.Create(r)
	if err != nil {
		logrus.WithFields(logrus.Fields{"reason": err}).Error("unable create bids record")
		return
	}

	r.ID = int(id)

	for _, car := range e.Claim.Cars {
		p := &model.Pass{
			CompanyBranch:     e.Claim.Company.Activity,
			CompanyOkved:      "",
			CompanyInn:        e.Claim.Company.INN,
			CompanyName:       e.Claim.Company.Title,
			CompanyAddress:    e.Claim.Company.Address,
			CompanyCeoPhone:   e.Claim.Company.Head.Contact.Phone,
			CompanyCeoEmail:   e.Claim.Company.Head.Contact.EMail,
			CompanyLastname:   e.Claim.Company.Head.FIO.Surname,
			CompanyFirstname:  e.Claim.Company.Head.FIO.Name,
			CompanyPatrname:   e.Claim.Company.Head.FIO.Patronymic,
			EmployeeLastname:  car.FIO.Surname,
			EmployeeFirstname: car.FIO.Name,
			EmployeePatrname:  car.FIO.Patronymic,
			EmployeeCar:       car.Number,
			EmployeeAgree:     1,
			EmployeeConfirm:   1,
			Source:            e.Source,
			District:          e.District,
			PassType:          e.PassType,
			PassNumber:        "",
			AlighnerPost:      "",
			AlighnerName:      "",
			SendType:          "",
			Status:            0,
			FileID:            e.FileID,
			CreatedAt:         time.Now(),
			CreatedBy:         e.CreatedBy,
			RequestID:         r.ID,
		}

		id, err := s.passRepo.Create(p)
		if err != nil {

			logrus.WithFields(logrus.Fields{"reason": err}).Error("unable create pass")

			s.UpdateState(&model.State{
				ID:     e.FileID,
				Status: 3,
				Error:  err,
			})

			return
		}
		p.ID = int(id)
	}

	f := &model.File{
		ID:        e.FileID,
		Status:    0,
		Log:       "",
		CreatedAt: time.Now(),
		Source:    e.Claim.Source,
	}

	if err := s.fileRepo.Update(f); err != nil {
		logrus.WithFields(logrus.Fields{"reason": err}).Error("unable update state")
	}
}

// UpdateState ...
func (s *EventService) UpdateState(e *model.State) {

	var log string
	if e.Error != nil {
		log = e.Error.Error()
	}

	f := &model.File{
		ID:        e.ID,
		Status:    e.Status,
		Log:       log,
		CreatedAt: time.Now(),
	}

	if err := s.fileRepo.Update(f); err != nil {
		logrus.WithFields(logrus.Fields{"reason": err}).Error("unable update state")
	}
}
