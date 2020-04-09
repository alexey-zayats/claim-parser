package services

import (
	"github.com/alexey-zayats/claim-parser/internal/model"
	"github.com/sirupsen/logrus"
	"go.uber.org/dig"
	"time"
)

// EventService ...
type EventService struct {
	bidSvc  *BidService
	fileSvc *FileService
	passSvc *PassService
}

// EventServiceDI ...
type EventServiceDI struct {
	dig.In
	BidSvc  *BidService
	FileSvc *FileService
	PassSvc *PassService
}

// NewEventService ...
func NewEventService(di EventServiceDI) *EventService {
	return &EventService{
		bidSvc:  di.BidSvc,
		fileSvc: di.FileSvc,
		passSvc: di.PassSvc,
	}
}

// StoreClaim ...
func (s *EventService) StoreClaim(claim *model.Claim) {

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
		logrus.WithFields(logrus.Fields{"reason": err}).Error("unable create bids record")
		return
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
		}

		if err := s.passSvc.Create(p); err != nil {

			logrus.WithFields(logrus.Fields{"reason": err}).Error("unable create pass")

			s.UpdateFile(event.FileID, 3, err.Error(), claim.Source)

			return
		}
	}

	s.UpdateFile(event.FileID, 0, "", claim.Source)
}

// UpdateFile ...
func (s *EventService) UpdateFile(id int64, status int, log string, source string) {

	f := &model.File{
		ID:        id,
		Status:    status,
		Log:       log,
		CreatedAt: time.Now(),
		Source:    source,
	}

	if err := s.fileSvc.UpdateState(f); err != nil {
		logrus.WithFields(logrus.Fields{"reason": err}).Error("unable update state")
	}
}
