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

// StoreClaim ...
func (s *EventService) StoreClaim(claim *model.Claim) {

	event := claim.Event

	r := &model.Request{
		FileID:         event.FileID,
		Status:         0,
		WorkflowStatus: 1,
		Code:           claim.Code,
		District:       event.District,
		PassType:       event.PassType,
		CreatedAt:      claim.Created,
		CreatedBy:      event.CreatedBy,
		UserID:         event.CreatedBy,
		Source:         claim.Source,
	}

	id, err := s.reqRepo.Create(r)
	if err != nil {
		logrus.WithFields(logrus.Fields{"reason": err}).Error("unable create bids record")
		return
	}

	r.ID = int(id)

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
			PassNumber:        "",
			AlighnerPost:      "",
			AlighnerName:      "",
			SendType:          "",
			Status:            0,
			FileID:            event.FileID,
			CreatedAt:         claim.Created,
			CreatedBy:         event.CreatedBy,
			RequestID:         r.ID,
		}

		id, err := s.passRepo.Create(p)
		if err != nil {

			logrus.WithFields(logrus.Fields{"reason": err}).Error("unable create pass")

			s.UpdateState(&model.State{
				ID:     event.FileID,
				Status: 3,
				Error:  err,
			})

			return
		}
		p.ID = int(id)
	}

	f := &model.File{
		ID:        event.FileID,
		Status:    0,
		Log:       "",
		CreatedAt: time.Now(),
		Source:    claim.Source,
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
