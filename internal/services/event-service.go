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
}

// EventServiceInput ...
type EventServiceInput struct {
	dig.In
	FileRepo interfaces.FileRepository
	PassRepo interfaces.PassRepository
}

// NewEventService ...
func NewEventService(param EventServiceInput) *EventService {
	return &EventService{
		fileRepo: param.FileRepo,
		passRepo: param.PassRepo,
	}
}

// StoreEvent ...
func (s *EventService) StoreEvent(e *model.Event) {

	for _, car := range e.Company.Cars {
		p := &model.Pass{
			CompanyBranch:     e.Company.Kind,
			CompanyOkved:      "",
			CompanyInn:        e.Company.INN,
			CompanyName:       e.Company.Name,
			CompanyAddress:    e.Company.Address,
			CompanyCeoPhone:   e.Company.Head.Contact.Phone,
			CompanyCeoEmail:   e.Company.Head.Contact.EMail,
			CompanyLastname:   e.Company.Head.FIO.Surname,
			CompanyFirstname:  e.Company.Head.FIO.Name,
			CompanyPatrname:   e.Company.Head.FIO.Patronymic,
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
		}

		if err := s.passRepo.Create(p); err != nil {

			s.UpdateState(&model.State{
				ID:     e.FileID,
				Status: 3,
				Error:  err,
			})

			logrus.WithFields(logrus.Fields{"reason": err}).Error("unable create pass")

			return
		}
	}

	s.UpdateState(&model.State{ID: e.FileID, Status: 0, Error: nil})
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
		logrus.WithFields(logrus.Fields{"reason": err}).Error("unable create pass")
	}
}
